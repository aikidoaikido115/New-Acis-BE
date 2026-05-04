-- Create "residents" table
CREATE TABLE "residents" (
  "id" text NOT NULL,
  "room_id" text NOT NULL,
  "first_name" text NOT NULL,
  "last_name" text NOT NULL,
  "age" bigint NOT NULL,
  "gender" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_residents_room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
