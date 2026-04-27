CREATE TABLE "support_tickets" (
  "id" text NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "subject" text NOT NULL,
  "message" text NOT NULL,
  "status" text NOT NULL DEFAULT 'open',
  "reporter_role" text NOT NULL,
  "created_by_user_id" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_support_tickets_status" CHECK ("status" IN ('open', 'in_progress', 'resolved')),
  CONSTRAINT "chk_support_tickets_reporter_role" CHECK ("reporter_role" IN ('Medical Staff', 'Kitchen Staff', 'Relative')),
  CONSTRAINT "fk_support_tickets_created_by_user" FOREIGN KEY ("created_by_user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX "idx_support_tickets_status" ON "support_tickets" ("status");
CREATE INDEX "idx_support_tickets_created_at" ON "support_tickets" ("created_at");
