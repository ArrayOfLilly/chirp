-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
    VALUES (
        gen_random_uuid(), 
        NOW(), 
        NOW(), 
        $1, 
        $2,
        false
        )
    RETURNING *;

-- name: GetUserByEmail :one
SELECT * 
    FROM users 
    WHERE email = $1;

-- name: GetUserByID :one
SELECT * 
    FROM users 
    WHERE id = $1;

-- name: UpdateUserData :one
UPDATE users 
    SET email = $2,
    hashed_password = $3,
    updated_at = NOW()
    WHERE id = $1
    RETURNING *;

-- name: UpgradeUserById :one
UPDATE users 
    SET is_chirpy_red = true,
    updated_at = NOW()
    WHERE id = $1
    RETURNING *;

-- name: DowngradeUserById :one
UPDATE users 
    SET is_chirpy_red = false,
    updated_at = NOW()
    WHERE id = $1
    RETURNING *;