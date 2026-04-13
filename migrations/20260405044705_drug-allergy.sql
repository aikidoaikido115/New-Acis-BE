-- Create "drug_allergies" table
CREATE TABLE "drug_allergies" (
  "id" text NOT NULL,
  "allergy_name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_drug_allergies_allergy_name" to table: "drug_allergies"
CREATE UNIQUE INDEX "idx_drug_allergies_allergy_name" ON "drug_allergies" ("allergy_name");
-- Create "resident_das" table
CREATE TABLE "resident_das" (
  "resident_id" text NOT NULL,
  "drug_allergy_id" text NOT NULL,
  "note_text" text NULL,
  "noted_at" timestamptz NOT NULL,
  PRIMARY KEY ("resident_id", "drug_allergy_id"),
  CONSTRAINT "fk_drug_allergies_resident_da" FOREIGN KEY ("drug_allergy_id") REFERENCES "drug_allergies" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_residents_resident_da" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
