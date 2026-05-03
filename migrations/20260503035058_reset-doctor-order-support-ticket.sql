-- Drop index "idx_drug_plans_pd_id_bkk_day_unique" from table: "drug_plans"
DROP INDEX "idx_drug_plans_pd_id_bkk_day_unique";
-- Create "doctor_orders" table
CREATE TABLE "doctor_orders" (
  "id" text NOT NULL,
  "resident_id" text NOT NULL,
  "order_date" text NULL,
  "order_type" text NULL,
  "title" text NOT NULL,
  "details" text NULL,
  "start_date" text NULL,
  "end_date" text NULL,
  "frequency" text NULL,
  "ordered_by" text NULL,
  "created_by_staff_id" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_doctor_orders_resident" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "support_tickets" table
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
  CONSTRAINT "fk_support_tickets_created_by_user" FOREIGN KEY ("created_by_user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
