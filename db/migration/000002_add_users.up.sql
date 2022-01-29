CREATE TABLE "users" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL ,
  "hashed_password" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "phone_number" varchar UNIQUE,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");

-- CREATE UNIQUE INDEX ON "accounts" ("owner", "property_id");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_property_id_key" UNIQUE ("owner","property_id");
