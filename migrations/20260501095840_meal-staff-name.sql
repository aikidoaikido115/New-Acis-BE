-- Add creator fields for meal history (safe for existing rows)
ALTER TABLE "meal_plans"
  ADD COLUMN IF NOT EXISTS "created_by_staff_id" text NULL,
  ADD COLUMN IF NOT EXISTS "staff_name" text NULL;
