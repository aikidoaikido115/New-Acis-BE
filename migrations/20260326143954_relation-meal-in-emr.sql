-- Create "meal_plans" table
CREATE TABLE "meal_plans" (
  "id" text NOT NULL,
  "menu_id" text NOT NULL,
  "back_up_menu_id" text NULL,
  "main_amount" smallint NOT NULL,
  "back_up_amount" smallint NULL,
  "is_allergy" boolean NULL,
  "meal_type" text NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_meal_plans_menu" FOREIGN KEY ("menu_id") REFERENCES "menus" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "allergies" table
CREATE TABLE "allergies" (
  "id" text NOT NULL,
  "allergy_name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_allergies_allergy_name" to table: "allergies"
CREATE UNIQUE INDEX "idx_allergies_allergy_name" ON "allergies" ("allergy_name");
-- Create "resident_allergies" table
CREATE TABLE "resident_allergies" (
  "resident_id" text NOT NULL,
  "allergy_id" text NOT NULL,
  "note_text" text NULL,
  "noted_at" timestamptz NOT NULL,
  PRIMARY KEY ("resident_id", "allergy_id"),
  CONSTRAINT "fk_allergies_resident_allergies" FOREIGN KEY ("allergy_id") REFERENCES "allergies" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_residents_resident_allergies" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
