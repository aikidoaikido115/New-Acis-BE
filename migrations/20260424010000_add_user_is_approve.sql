-- Add approval flag to users
ALTER TABLE "users"
ADD COLUMN "is_approve" boolean NOT NULL DEFAULT false;

-- Preserve existing users so the deployment does not lock out current accounts
UPDATE "users"
SET "is_approve" = true;