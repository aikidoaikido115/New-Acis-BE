-- Create "rooms" table
CREATE TABLE "rooms" (
  "id" text NOT NULL,
  "staff_id" text NULL,
  "floor" smallint NOT NULL,
  "room_number" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_rooms_staff" FOREIGN KEY ("staff_id") REFERENCES "staffs" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_rooms_room_number" to table: "rooms"
CREATE UNIQUE INDEX "idx_rooms_room_number" ON "rooms" ("room_number");
