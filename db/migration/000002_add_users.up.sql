CREATE TABLE "users" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL ,
  "hashed_password" varchar NOT NULL,
  "nickname" varchar NOT NULL,
  "phone_number" varchar UNIQUE,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

-- CREATE UNIQUE INDEX ON "accounts" ("user_id", "property_id");
ALTER TABLE "accounts" ADD CONSTRAINT "user_id_property_id_key" UNIQUE ("user_id","property_id");
