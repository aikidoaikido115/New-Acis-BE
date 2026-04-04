-- Create "drug_masters" table
CREATE TABLE "drug_masters" (
  "id" text NOT NULL,
  "name" text NOT NULL,
  "dose" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_drug_masters_name" UNIQUE ("name")
);
