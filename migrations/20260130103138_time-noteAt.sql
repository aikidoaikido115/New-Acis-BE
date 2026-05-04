-- Modify "resident_labels" table
ALTER TABLE "resident_labels" ALTER COLUMN "noted_at" TYPE timestamptz USING noted_at::timestamptz;
