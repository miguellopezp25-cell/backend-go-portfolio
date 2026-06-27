package main

import (
	"github.com/miguel/go-back-portfolo/cmd"
)

// @title Visitor Management API
// @version 1.0.0
// @description REST API for managing visitor registrations with PostgreSQL and JSONB storage.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cmd.Execute()
}
