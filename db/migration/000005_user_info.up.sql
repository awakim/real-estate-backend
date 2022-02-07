CREATE TABLE "user_information" (
  "user_id" uuid PRIMARY KEY,
  "firstname" varchar NOT NULL,
  "lastname" varchar NOT NULL,
  "phone_number" varchar UNIQUE NOT NULL,
  "nationality" varchar NOT NULL,
  "address" varchar NOT NULL,
  "postal_code" varchar NOT NULL,
  "city" varchar NOT NULL,
  "country" varchar NOT NULL
);

COMMENT ON COLUMN "user_information"."phone_number" IS 'has to be E164 compliant';

ALTER TABLE "user_information" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");