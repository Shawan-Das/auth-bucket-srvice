-- name: CreateUser :one
INSERT INTO common.users (
    user_name,
    email,
    phone,
    pass,
    role,
    pss_valid,
    otp_valid
) VALUES (
    $1, $2, $3, $4, $5, true, false
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM common.users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM common.users
WHERE user_id = $1 LIMIT 1;

-- name: UpdateRefreshToken :exec
UPDATE common.users
SET refresh_token = $1,
    refresh_token_exp = $2,
    updated_at = now()
WHERE user_id = $3;

-- name: GetUserByRefreshToken :one
SELECT * FROM common.users
WHERE refresh_token = $1
AND refresh_token_exp > now()
LIMIT 1;

-- name: InvalidateRefreshToken :exec
UPDATE common.users
SET refresh_token = '',
    refresh_token_exp = NULL,
    updated_at = now()
WHERE user_id = $1;

-- name: UpdateUserPassword :exec
UPDATE common.users
SET pass = $1,
    pss_valid = true,
    updated_at = now()
WHERE user_id = $2;

-- name: GetUserByEmailWithRole :one
SELECT * FROM common.users
WHERE email = $1 AND role = $2 LIMIT 1;
