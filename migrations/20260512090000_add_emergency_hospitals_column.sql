-- Add emergency_hospitals column to residents for multi-hospital support
ALTER TABLE "residents"
  ADD COLUMN IF NOT EXISTS "emergency_hospitals" jsonb;
