// Package model define las estructuras de dominio del proyecto. Están en schema/
// porque son la representación en código del esquema de datos, independientemente
// de cómo se almacenen (PostgreSQL JSONB, memoria, etc.).
package schema

// Visitor es el modelo completo de dominio para un visitante.
type Visitor struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Country string `json:"country"`
	City    string `json:"city"`
}
