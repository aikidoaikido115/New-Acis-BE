-- Create "nurse_notes" table
CREATE TABLE "nurse_notes" (
  "id" text NOT NULL,
  "resident_id" text NOT NULL,
  "category" text NOT NULL,
  "content" text NOT NULL,
  "priority" text NOT NULL,
  "send_note" boolean NOT NULL DEFAULT false,
  "created_by_staff_id" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_nurse_notes_resident" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create "wound_care_notes" table
CREATE TABLE "wound_care_notes" (
  "id" text NOT NULL,
  "resident_id" text NOT NULL,
  "location" text NOT NULL,
  "wound_type" text NOT NULL,
  "size" text NULL,
  "treatment" text NULL,
  "supplies" text NULL,
  "status" text NULL,
  "image_url" text NULL,
  "note" text NULL,
  "created_by_staff_id" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_wound_care_notes_resident" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create "relative_notes" table
CREATE TABLE "relative_notes" (
  "id" text NOT NULL,
  "resident_id" text NOT NULL,
  "relation" text NOT NULL,
  "content" text NOT NULL,
  "send_note" boolean NOT NULL DEFAULT true,
  "created_by_staff_id" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_relative_notes_resident" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
