ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "user_id_property_id_key";

ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_user_id_fkey";

DROP TABLE IF EXISTS "users";