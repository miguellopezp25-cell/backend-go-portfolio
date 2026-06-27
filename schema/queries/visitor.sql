-- name: CreateVisitor :one
INSERT INTO visitor.visitors ("data")
VALUES ($1)
RETURNING *;

-- name: GetVisitor :one
SELECT id, data, created_at
FROM visitor.visitors
WHERE id = $1 LIMIT 1;

-- name: ListVisitors :many
SELECT id, data, created_at
FROM visitor.visitors
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountVisitors :one
SELECT COUNT(*) FROM visitor.visitors;

-- name: UpdateVisitor :one
UPDATE visitor.visitors
SET data = $1
WHERE id = $2
RETURNING *;

-- name: DeleteVisitor :one
DELETE FROM visitor.visitors
WHERE id = $1
RETURNING id;


