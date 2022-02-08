-- name: CreateUserInfo :one
INSERT INTO user_information (
  user_id,
  firstname,
  lastname,
  phone_number,
  nationality,
  address,
  postal_code,
  city,
  country
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;


-- name: GetUserInfo :one
SELECT * FROM user_information
WHERE user_id = $1 LIMIT 1;