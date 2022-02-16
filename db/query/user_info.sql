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
) ON CONFLICT (user_id) DO UPDATE
  SET 
  firstname = excluded.firstname,
  lastname = excluded.lastname,
  phone_number = excluded.phone_number,
  nationality = excluded.nationality,
  address = excluded.address,
  postal_code = excluded.postal_code,
  city = excluded.city,
  country = excluded.country 
RETURNING *;


-- name: GetUserInfo :one
SELECT * FROM user_information
WHERE user_id = $1 LIMIT 1;

-- name: ExistsUserInfo :one
SELECT EXISTS(
  SELECT 1 FROM user_information
  WHERE user_id = $1 LIMIT 1
);