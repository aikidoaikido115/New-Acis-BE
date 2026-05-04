-- Create "audit_logs" table
CREATE TABLE "audit_logs" (
  "id" text NOT NULL,
  "table_name" text NOT NULL,
  "record_id" text NOT NULL,
  "user_id" text NOT NULL,
  "action" text NOT NULL,
  "old_value" text NULL,
  "new_value" text NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_audit_logs_table_name" UNIQUE ("table_name")
);
-- Create "roles" table
CREATE TABLE "roles" (
  "id" text NOT NULL,
  "name" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_roles_name" UNIQUE ("name")
);
-- Create "users" table
CREATE TABLE "users" (
  "id" text NOT NULL,
  "role_id" text NOT NULL,
  "username" text NOT NULL,
  "number_of_usernames" bigint NULL DEFAULT 0,
  "email" text NOT NULL,
  "password" text NULL,
  "first_name" text NULL,
  "last_name" text NULL,
  "nickname" text NULL,
  "gender" text NULL,
  "profile_image" text NULL DEFAULT 'https://www.isranews.org/article/images/2025/Harry/6/Hun_Sen_July_2019.jpg',
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email"),
  CONSTRAINT "uni_users_username" UNIQUE ("username"),
  CONSTRAINT "fk_users_role" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "otps" table
CREATE TABLE "otps" (
  "user_id" text NOT NULL,
  "otp" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  PRIMARY KEY ("user_id"),
  CONSTRAINT "uni_otps_otp" UNIQUE ("otp"),
  CONSTRAINT "fk_otps_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "staffs" table
CREATE TABLE "staffs" (
  "id" text NOT NULL,
  "user_id" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_staffs_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "staffs_files" table
CREATE TABLE "staffs_files" (
  "id" text NOT NULL,
  "staff_id" text NOT NULL,
  "file_name" text NOT NULL,
  "file_type" text NOT NULL,
  "file_size" bigint NOT NULL,
  "file" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_staffs_files_staff" FOREIGN KEY ("staff_id") REFERENCES "staffs" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "temp_tokens" table
CREATE TABLE "temp_tokens" (
  "user_id" text NOT NULL,
  "token" text NOT NULL,
  PRIMARY KEY ("user_id"),
  CONSTRAINT "uni_temp_tokens_token" UNIQUE ("token"),
  CONSTRAINT "fk_temp_tokens_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
