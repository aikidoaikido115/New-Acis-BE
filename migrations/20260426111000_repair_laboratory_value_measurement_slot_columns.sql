-- Repair migration for environments that ran the old laboratory slot migration containing rollback SQL.
-- This migration is idempotent and safe to run multiple times.
ALTER TABLE "laboratory_values"
ADD COLUMN IF NOT EXISTS "measurement_date" date;

UPDATE "laboratory_values"
SET "measurement_date" = ("created_at" AT TIME ZONE 'Asia/Bangkok')::date
WHERE "measurement_date" IS NULL;

ALTER TABLE "laboratory_values"
ALTER COLUMN "measurement_date" SET NOT NULL;

ALTER TABLE "laboratory_values"
ADD COLUMN IF NOT EXISTS "time_of_day" text;

UPDATE "laboratory_values"
SET "time_of_day" = 'เช้า'
WHERE "time_of_day" IS NULL OR BTRIM("time_of_day") = '';

ALTER TABLE "laboratory_values"
ALTER COLUMN "time_of_day" SET NOT NULL;

CREATE INDEX IF NOT EXISTS "idx_laboratory_values_measurement_date"
ON "laboratory_values" ("measurement_date");

CREATE INDEX IF NOT EXISTS "idx_laboratory_values_resident_measurement_date"
ON "laboratory_values" ("resident_id", "measurement_date");
