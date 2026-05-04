DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = 'public'
      AND table_name = 'support_tickets'
  ) THEN
    ALTER TABLE "support_tickets"
      DROP CONSTRAINT IF EXISTS "chk_support_tickets_reporter_role";

    ALTER TABLE "support_tickets"
      ADD CONSTRAINT "chk_support_tickets_reporter_role"
      CHECK ("reporter_role" IN ('Medical Staff', 'Kitchen Staff', 'Relative', 'Super User', 'Admin'));
  END IF;
END $$;