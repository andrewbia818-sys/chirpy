-- name: DeleteUsers :one
DELETE FROM users
RETURNING *;