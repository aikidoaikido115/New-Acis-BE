-- Modify "residents" table
ALTER TABLE "residents" DROP COLUMN "adl_assessment";
-- Create "laboratory_values" table
CREATE TABLE "laboratory_values" (
  "id" text NOT NULL,
  "resident_id" text NOT NULL,
  "blood_glucose" numeric NULL,
  "fluid_in" numeric NULL,
  "fluid_out" numeric NULL,
  "urine_output" numeric NULL,
  "urine_type" text NULL,
  "stool" smallint NULL,
  "diaper_change" smallint NULL,
  "created_by_staff_id" text NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id", "resident_id"),
  CONSTRAINT "fk_laboratory_values_resident" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
