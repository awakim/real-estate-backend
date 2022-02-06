-- name: CreateUserInfo :one
INSERT INTO user_information (
  user_id,
  firstname,
  lastname,
  phone_number,
  nationality,
  gender,
  address,
  postal_code,
  city,
  country
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) ON CONFLICT (user_id) DO UPDATE
  SET 
  firstname = excluded.firstname,
  lastname = excluded.lastname,
  phone_number = excluded.phone_number,
  nationality = excluded.nationality,
  gender = excluded.gender,
  address = excluded.address,
  postal_code = excluded.postal_code,
  city = excluded.city,
  country = excluded.country 
RETURNING *;


-- name: GetUserInfo :one
SELECT * FROM user_information
WHERE user_id = $1 LIMIT 1;