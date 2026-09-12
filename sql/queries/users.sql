-- name: CreateUser :one
INSERT INTO users (
    email,
    hashed_password,
    is_chirpy_red,
    created_at,
    updated_at
)
VALUES (
    $1,
    $2,
    FALSE,
    NOW(),
    NOW()
)
RETURNING id, email, created_at, updated_at, is_chirpy_red;
