-- Create "activities" table
CREATE TABLE "activities" (
  "id" text NOT NULL,
  "staff_id" text NOT NULL,
  "activity_name" text NOT NULL,
  "activity_type" text NOT NULL,
  "description" text NULL,
  "location" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_activities_staff" FOREIGN KEY ("staff_id") REFERENCES "staffs" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
