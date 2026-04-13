-- Create "personal_drugs" table
CREATE TABLE "personal_drugs" (
  "id" text NOT NULL,
  "resident_id" text NOT NULL,
  "dm_id" text NOT NULL,
  "amount" text NOT NULL,
  "amount_unit" text NOT NULL,
  "frequency" bigint NOT NULL,
  "time_of_day" text NOT NULL,
  "timing" text NOT NULL,
  "description" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_personal_drugs_drug_master" FOREIGN KEY ("dm_id") REFERENCES "drug_masters" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_personal_drugs_resident" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
