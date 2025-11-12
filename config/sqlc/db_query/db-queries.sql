-- --------------------- AUTHENTICATION ------------------------------
-- name: GetUserByEmail :one
SELECT user_id, user_name, email, phone, pass, pss_valid, otp, otp_valid, otp_exp, role 
FROM common.users 
WHERE email = $1;

-- name: GetUserByLogin :one
SELECT user_id, user_name, email, phone, pass, pss_valid, otp, otp_valid, otp_exp, role 
FROM common.users 
WHERE user_name = $1 OR email = $1 OR phone = $1;

-- name: CreateUser :exec
INSERT INTO common.users(user_name, email, phone, pass, role) 
VALUES($1, $2, $3, $4, $5);

-- name: UpdatePassword :exec
UPDATE common.users 
SET pass = $1, pss_valid = $2 
WHERE email = $3;

-- name: GetAllUsers :many
SELECT user_id, user_name, email 
FROM common.users 
ORDER BY user_id;

-- --------------------- SATCOM DATA ------------------------------
-- name: CreateSatcomData :exec
INSERT INTO common.satcom_data(company, category, "type", "date", "time", db_port, ui_port, url, ip, status)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetSatcomDataById :one
SELECT id, company, category, "type", "date", "time", db_port, ui_port, url, ip, status
FROM common.satcom_data
WHERE id = $1;

-- name: GetAllSatcomData :many
SELECT id, company, category, "type", "date", "time", db_port, ui_port, url, ip, status
FROM common.satcom_data
ORDER BY id;

-- name: UpdateSatcomData :exec
UPDATE common.satcom_data
SET company = $1, category = $2, "type" = $3, "date" = $4, "time" = $5, 
    db_port = $6, ui_port = $7, url = $8, ip = $9, status = $10
WHERE id = $11;

-- name: DeleteSatcomData :exec
DELETE FROM common.satcom_data
WHERE id = $1;