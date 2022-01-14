-- name: CreateProperty :one
INSERT INTO properties (
  "name",
  "description",
  initial_block_count,
  remaining_block_count
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetProperty :one
SELECT * FROM properties
WHERE id = $1 LIMIT 1;