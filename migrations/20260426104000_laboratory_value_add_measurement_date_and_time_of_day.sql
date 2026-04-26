-- +goose Up
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

-- +goose Down
DROP INDEX IF EXISTS "idx_laboratory_values_resident_measurement_date";
DROP INDEX IF EXISTS "idx_laboratory_values_measurement_date";

ALTER TABLE "laboratory_values"
DROP COLUMN IF EXISTS "time_of_day";

ALTER TABLE "laboratory_values"
DROP COLUMN IF EXISTS "measurement_date";
