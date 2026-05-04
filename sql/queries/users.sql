-- name: CreateUsers :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
   $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :many
DELETE FROM users 
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
