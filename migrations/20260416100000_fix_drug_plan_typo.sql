-- Rename a column from "is_ommitted" to "is_omitted"
ALTER TABLE "drug_plans" RENAME COLUMN "is_ommitted" TO "is_omitted";
-- Rename a column from "ommitted_reason" to "omitted_reason"
ALTER TABLE "drug_plans" RENAME COLUMN "ommitted_reason" TO "omitted_reason";
