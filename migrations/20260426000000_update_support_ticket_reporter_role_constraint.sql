ALTER TABLE "support_tickets" DROP CONSTRAINT "chk_support_tickets_reporter_role";

ALTER TABLE "support_tickets"
  ADD CONSTRAINT "chk_support_tickets_reporter_role" CHECK ("reporter_role" IN ('Medical Staff', 'Kitchen Staff', 'Relative', 'Super User', 'Admin'));