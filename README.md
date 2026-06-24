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

## Docker

### docker-compose.yml

```yaml
services:
  app:                          # Servicio de la app Go
    build:
      context: .                # Usa el Dockerfile del proyecto
      dockerfile: Dockerfile
    ports:
      - "8080:8080"             # Expone puerto 8080 al host
    depends_on:
      db:
        condition: service_healthy  # Espera a que PostgreSQL esté listo
    environment:
      - CONFIG_PATH=/app/config.yaml   # Ruta del config dentro del contenedor
      - DATABASE_HOST=db               # Apunta a la BD por nombre del servicio

  db:                           # Servicio PostgreSQL
    image: postgres:16-alpine
    ports:
      - "5432:5432"             # Expone PostgreSQL al host (para DBeaver)
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: visitors
    volumes:
      - pgdata:/var/lib/postgresql/data   # Persiste datos aunque mate el contenedor
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pgdata:                       # Volumen nombrado para persistencia
```

### Dockerfile

```dockerfile
# === ETAPA 1: compilación ===
FROM golang:1.24-alpine AS builder
WORKDIR /app                      # Directorio de trabajo
COPY go.mod go.sum ./             # Solo mod/sum primero (cachea dependencias)
RUN go mod download               # Descarga dependencias
COPY . .                          # Copia todo el código
RUN CGO_ENABLED=0 go build -o bin/app ./main.go  # Compila binario estático

# === ETAPA 2: imagen final mínima ===
FROM alpine:3.21
RUN apk add --no-cache ca-certificates   # Certificados TLS
WORKDIR /app
COPY --from=builder /app/bin/app .              # Solo el binario
COPY --from=builder /app/config.yaml .          # Config
COPY --from=builder /app/schema/migrations ./schema/migrations  # Migraciones SQL
COPY --from=builder /app/entrypoint.sh .        # Script de inicio
RUN chmod +x entrypoint.sh
EXPOSE 8080                        # Puerto que usa la app
CMD ["./entrypoint.sh"]            # Entrypoint: corre migrate luego serve
```

### entrypoint.sh

```sh
#!/bin/sh
set -e                    # Detiene si algún comando falla
./app migrate            # Crea tablas automáticamente
exec ./app serve         # Reemplaza el shell por el proceso del server
```

## Comandos disponibles

| Comando                   | Descripción                    |
|---------------------------|--------------------------------|
| `go run main.go serve`    | Iniciar servidor HTTP          |
| `go run main.go migrate`  | Ejecutar migraciones           |
| `go run main.go version`  | Mostrar versión                |
