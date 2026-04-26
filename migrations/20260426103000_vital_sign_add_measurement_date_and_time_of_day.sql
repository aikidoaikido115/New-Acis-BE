-- Add per-day slot tracking fields for vital signs.
ALTER TABLE "vital_signs"
ADD COLUMN IF NOT EXISTS "measurement_date" date;

UPDATE "vital_signs"
SET "measurement_date" = ("created_at" AT TIME ZONE 'Asia/Bangkok')::date
WHERE "measurement_date" IS NULL;

ALTER TABLE "vital_signs"
ALTER COLUMN "measurement_date" SET NOT NULL;

ALTER TABLE "vital_signs"
ADD COLUMN IF NOT EXISTS "time_of_day" text;

UPDATE "vital_signs"
SET "time_of_day" = 'เช้า'
WHERE "time_of_day" IS NULL OR BTRIM("time_of_day") = '';

ALTER TABLE "vital_signs"
ALTER COLUMN "time_of_day" SET NOT NULL;

CREATE INDEX IF NOT EXISTS "idx_vital_signs_measurement_date"
ON "vital_signs" ("measurement_date");

CREATE INDEX IF NOT EXISTS "idx_vital_signs_resident_measurement_date"
ON "vital_signs" ("resident_id", "measurement_date");
