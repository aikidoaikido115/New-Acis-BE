-- Modify "drug_masters" table
ALTER TABLE "drug_masters" DROP CONSTRAINT "uni_drug_masters_name";
-- Create index "idx_drug_masters_name_dose" to table: "drug_masters"
CREATE UNIQUE INDEX "idx_drug_masters_name_dose" ON "drug_masters" ("name", "dose");
