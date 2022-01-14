CREATE TABLE "properties" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" text NOT NULL,
  "initial_block_count" bigint NOT NULL,
  "remaining_block_count" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("property_id") REFERENCES "properties" ("id");

COMMENT ON COLUMN "properties"."initial_block_count" IS 'must be greater than or equal to zero';

COMMENT ON COLUMN "properties"."remaining_block_count" IS 'must be greater than or equal to zero';
