-- name: GetNote :one
SELECT * FROM notes
WHERE id = $1 LIMIT 1;

-- name: ListNotesByUserID :many
SELECT * FROM notes WHERE user_id = $1;

-- name: CreateNote :one
INSERT INTO notes (
  title, content, user_id, created_at
) VALUES (
  $1, $2, $3, NOW()
)
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = $1 AND user_id = $2;

-- name: UpdateNote :exec
UPDATE notes
SET title = $1, content = $2
WHERE id = $3 AND user_id = $4;

