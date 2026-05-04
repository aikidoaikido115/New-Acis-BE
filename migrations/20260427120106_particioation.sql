-- Create "participations" table
CREATE TABLE "participations" (
  "resident_id" text NOT NULL,
  "as_id" text NOT NULL,
  "is_participating" boolean NOT NULL,
  "img_urls" jsonb NOT NULL DEFAULT '[]',
  PRIMARY KEY ("resident_id", "as_id"),
  CONSTRAINT "fk_participations_activity_schedule" FOREIGN KEY ("as_id") REFERENCES "activity_schedules" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_participations_resident" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
