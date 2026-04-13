-- Create "drug_plans" table
CREATE TABLE "drug_plans" (
  "id" text NOT NULL,
  "pd_id" text NOT NULL,
  "is_taken" boolean NOT NULL,
  "taken_at" timestamptz NULL,
  "given_by_staff_id" text NOT NULL,
  "is_ommitted" boolean NOT NULL,
  "ommitted_reason" text NULL,
  "notes" text NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_drug_plans_personal_drug" FOREIGN KEY ("pd_id") REFERENCES "personal_drugs" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
