ALTER TABLE "users"
  ADD COLUMN IF NOT EXISTS "phone" text;

ALTER TABLE "residents"
  ALTER COLUMN "id_card_number" DROP NOT NULL,
  ALTER COLUMN "check_in_date" DROP NOT NULL,
  ALTER COLUMN "room_id" DROP NOT NULL;

ALTER TABLE "residents"
  ADD COLUMN IF NOT EXISTS "emergency_contacts" jsonb,
  ADD COLUMN IF NOT EXISTS "profile_image" text;