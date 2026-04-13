-- Modify "personal_drugs" table
ALTER TABLE "personal_drugs" ADD COLUMN "created_at" timestamptz NOT NULL, ADD COLUMN "updated_at" timestamptz NOT NULL;
