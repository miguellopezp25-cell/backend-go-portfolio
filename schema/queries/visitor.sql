-- name: CreateVisitor :one
INSERT INTO visitor.visitors ("data")
VALUES ($1)
RETURNING *;

-- name: GetVisitor :one
SELECT id, data, created_at
FROM visitor.visitors
WHERE id = $1 LIMIT 1;


