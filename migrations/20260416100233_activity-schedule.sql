-- Create "activity_schedules" table
CREATE TABLE "activity_schedules" (
  "id" text NOT NULL,
  "activity_id" text NOT NULL,
  "date" timestamptz NOT NULL,
  "start_time" timestamptz NOT NULL,
  "end_time" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_activity_schedules_activity" FOREIGN KEY ("activity_id") REFERENCES "activities" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
