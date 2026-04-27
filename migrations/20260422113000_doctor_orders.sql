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
  "image_url" text NULL,
  "created_by_staff_id" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_doctor_orders_resident" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
