-- name: CreateSession :one
INSERT INTO sessions (
  id, user_id, refresh_token, is_revoked, expires_at, created_at
) VALUES (
  $1, $2, $3, $4, $5, NOW()
)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: RevokeSession :exec
UPDATE sessions
SET is_revoked = TRUE
WHERE id = $1;

-- name: RevokeUserSessions :exec
UPDATE sessions
SET is_revoked = TRUE
WHERE user_id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;


