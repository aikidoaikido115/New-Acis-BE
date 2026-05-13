-- Add series_id and status to activity_schedules
ALTER TABLE "activity_schedules"
  ADD COLUMN IF NOT EXISTS "series_id" text,
  ADD COLUMN IF NOT EXISTS "status" text NOT NULL DEFAULT 'active';

CREATE INDEX IF NOT EXISTS "idx_activity_schedules_series_id" ON "activity_schedules" ("series_id");
CREATE INDEX IF NOT EXISTS "idx_activity_schedules_status" ON "activity_schedules" ("status");
