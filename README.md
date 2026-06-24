# backend-go-portfolio

API REST en Go para portfolio. PostgreSQL con JSONB, sqlc, Gin, arquitectura por capas.

## Requisitos

- Go 1.24+
- Docker y Docker Compose

## Pasos (desarrollo local)

```bash
# 1. Levantar PostgreSQL
docker compose up -d db

# 2. Crear tablas
go run main.go migrate

# 3. Iniciar servidor
go run main.go serve
```

## Pasos (Docker completo)

```bash
# 1. Construir imagen y levantar todo (app + db)
docker compose up -d --build

# 2. Crear tablas dentro del contenedor
docker compose exec app ./app migrate

# 3. Ver logs
docker compose logs -f
```

Servidor en `http://localhost:8080`.

## Endpoints

| Método | Ruta                  | Descripción           |
|--------|-----------------------|-----------------------|
| GET    | `/healthz`            | Health check          |
| POST   | `/api/v1/visitors`    | Crear visitor         |
| GET    | `/api/v1/visitors/:id`| Obtener visitor por ID|

### Ejemplo

```bash
curl -X POST http://localhost:8080/api/v1/visitors \
  -H "Content-Type: application/json" \
  -d '{"name":"Miguel","email":"m@test.com","country":"MX","city":"CDMX"}'
```

## Estructura del proyecto

```
main.go              ← Entry point
cmd/                 ← CLI commands (Cobra)
├── serve.go         ←   Inicia servidor HTTP
├── migrate.go       ←   Ejecuta migraciones SQL
├── version.go       ←   Muestra versión
└── root.go          ←   Raíz de comandos

api/                 ← Capa HTTP (Gin)
├── server.go        ←   Server: conecta DB, crea store/service, inicia HTTP
├── router.go        ←   Registra rutas en Gin
├── health.go        ←   GET /healthz
├── visitor_handler.go  ←   POST /api/v1/visitors
└── getvisitor_handler.go ← GET /api/v1/visitors/:id

service/visitorservice/  ← Lógica de negocio
├── service.go       ←   Service struct + VisitorRequest
├── create.go        ←   Crear visitor
├── get.go           ←   Obtener visitor por ID
└── service_test.go  ←   Tests unitarios

schema/              ← Modelos de dominio
├── visitor.go       ←   Visitor struct (dominio)
├── queries/         ←   SQL queries para sqlc
└── db/              ←   Código generado por sqlc
    ├── db.go        ←   DBTX interface, Queries
    ├── querier.go   ←   Querier interface
    ├── store.go     ←   Store: embed Querier + ExecTx
    ├── models.go    ←   VisitorVisitor (DB model)
    └── visitor.sql.go ←   Implementación de queries

config/              ← Configuración (Viper + YAML)
├── config.go        ←   Carga config.yaml con vars de entorno

database/            ← Conexión a BD
├── database.go      ←   Pool de conexiones pgxpool

pkg/                 ← Utilidades compartidas
├── response/        ←   Envelope JSON estándar (success, data, error)
└── errors/          ←   Errores personalizados (ErrNotFound)
```

## Capas (flujo de una petición)

```
Cliente HTTP
    ↓
api/ (Gin handler)     ← valida request, llama al service
    ↓
service/               ← lógica de negocio, llama al store
    ↓
schema/db/ (Store)     ← embed Querier, executa queries SQL
    ↓
PostgreSQL             ← datos en JSONB
```

## Comandos disponibles

| Comando                   | Descripción                    |
|---------------------------|--------------------------------|
| `go run main.go serve`    | Iniciar servidor HTTP          |
| `go run main.go migrate`  | Ejecutar migraciones           |
| `go run main.go version`  | Mostrar versión                |
