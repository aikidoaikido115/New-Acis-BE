-- Create "relatives" table
CREATE TABLE "relatives" (
  "id" text NOT NULL,
  "user_id" text NOT NULL,
  "resident_id" text NOT NULL,
  "relative_password" text NOT NULL,
  "relation" text NOT NULL,
  "phone" text NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_relatives_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_relatives_resident" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE UNIQUE INDEX "uidx_relatives_resident_id" ON "relatives" ("resident_id");
CREATE UNIQUE INDEX "uidx_relatives_user_id" ON "relatives" ("user_id");

-- Create "relative_magic_link_tokens" table
CREATE TABLE "relative_magic_link_tokens" (
  "id" text NOT NULL,
  "relative_id" text NOT NULL,
  "resident_id" text NOT NULL,
  "token" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "last_accessed_at" timestamptz NULL,
  "created_by_user_id" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_relative_magic_link_tokens_relative" FOREIGN KEY ("relative_id") REFERENCES "relatives" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_relative_magic_link_tokens_resident" FOREIGN KEY ("resident_id") REFERENCES "residents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_relative_magic_link_tokens_created_by_user" FOREIGN KEY ("created_by_user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE UNIQUE INDEX "uidx_relative_magic_link_tokens_token" ON "relative_magic_link_tokens" ("token");
CREATE INDEX "idx_relative_magic_link_tokens_resident_id" ON "relative_magic_link_tokens" ("resident_id");
CREATE INDEX "idx_relative_magic_link_tokens_expires_at" ON "relative_magic_link_tokens" ("expires_at");
