ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_property_id_key";

ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

DROP TABLE IF EXISTS "users";