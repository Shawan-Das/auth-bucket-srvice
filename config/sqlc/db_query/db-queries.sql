-- --------------------- AUTHENTICATION ------------------------------
-- name: GetUserByEmail :one
SELECT user_id, user_name, email, phone, pass, pss_valid, otp, otp_valid, otp_exp, role 
FROM common.users 
WHERE email = $1;

-- name: CreateUser :exec
INSERT INTO common.users(user_name, email, phone, pass, role) 
VALUES($1, $2, $3, $4, $5);

-- name: UpdatePassword :exec
UPDATE common.users 
SET pass = $1, pss_valid = $2 
WHERE email = $3;
