-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = $2,
    updated_at = $3
WHERE token = $1
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at;