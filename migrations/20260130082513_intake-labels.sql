-- Create "intake_labels" table
CREATE TABLE "intake_labels" (
  "id" text NOT NULL,
  "label_name" text NOT NULL,
  "description" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "resident_labels" table
CREATE TABLE "resident_labels" (
  "resident_id" text NOT NULL,
  "label_id" text NOT NULL,
  "note_text" text NOT NULL,
  "noted_at" text NOT NULL,
  PRIMARY KEY ("resident_id", "label_id"),
  CONSTRAINT "fk_intake_labels_resident_labels" FOREIGN KEY ("label_id") REFERENCES "intake_labels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_residents_resident_labels" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
