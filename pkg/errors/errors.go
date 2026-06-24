// Package errors define errores del dominio que pueden atravesar las capas
// (handler → service → repository) sin acoplar el código a errores de
// librerías externas como pgx.
package errors

import "errors"

var (
	ErrNotFound = errors.New("resource not found")
)
