-- Create "vital_signs" table
CREATE TABLE "vital_signs" (
  "id" text NOT NULL,
  "resident_id" text NOT NULL,
  "temperature" text NULL,
  "heart_rate" text NULL,
  "breathing_rate" text NULL,
  "blood_pressure_systolic" text NULL,
  "blood_pressure_diastolic" text NULL,
  "oxygen_saturation" text NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id", "resident_id"),
  CONSTRAINT "fk_vital_signs_resident" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
