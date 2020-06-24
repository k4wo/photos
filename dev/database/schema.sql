DROP TYPE IF EXISTS "public"."file_type";
CREATE TYPE "public"."file_type" AS ENUM ('IMAGE', 'VIDEO', 'ANIMATION', 'COLLAGE');
-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS users_id_seq;
-- Table Definition
CREATE TABLE IF NOT EXISTS "public"."users" (
  "id" int4 NOT NULL DEFAULT nextval('users_id_seq' :: regclass),
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  PRIMARY KEY ("id")
);
-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS files_id_seq;
-- Table Definition
CREATE TABLE IF NOT EXISTS "public"."files" (
  "id" int4 NOT NULL DEFAULT nextval('files_id_seq' :: regclass),
  "type" "public"."file_type" NOT NULL,
  "owner" int4 NOT NULL,
  "name" varchar NOT NULL,
  "hash" varchar NOT NULL,
  "size" int4 NOT NULL DEFAULT '0' :: bigint,
  "extension" varchar NOT NULL,
  "mime" varchar NOT NULL,
  "latitude" float8,
  "longitude" float8,
  "orientation" int2,
  "model" varchar,
  "camera" varchar,
  "iso" int2,
  "focal_length" float4,
  "exposure_time" varchar,
  "f_number" decimal(5, 1),
  "width" int2,
  "height" int2,
  "date" timestamptz,
  "updated_at" timestamptz DEFAULT now(),
  "created_at" timestamptz DEFAULT now(),
  CONSTRAINT "files_owner_fkey" FOREIGN KEY ("owner") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  PRIMARY KEY ("id")
);
-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS albums_id_seq;
-- Table Definition
CREATE TABLE "public"."albums" (
    "id" int4 NOT NULL DEFAULT nextval('albums_id_seq'::regclass),
    "owner" int4 NOT NULL,
    "name" varchar NOT NULL,
    "size" int4 NOT NULL DEFAULT '0'::bigint,
    "cover" int4,
    "updated_at" time DEFAULT now(),
    "created_at" time DEFAULT now(),
    CONSTRAINT "albums_cover_fkey" FOREIGN KEY ("cover") REFERENCES "public"."files"("id") ON DELETE CASCADE,
    CONSTRAINT "albums_owner_fkey" FOREIGN KEY ("owner") REFERENCES "public"."users"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);
-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS album_file_id_seq;
-- Table Definition
CREATE TABLE  IF NOT EXISTS "public"."album_file" (
    "id" int4 NOT NULL DEFAULT nextval('album_file_id_seq'::regclass),
    "album" int4 NOT NULL,
    "added_by" int4,
    "file" int4 NOT NULL,
    "updated_at" timestamptz DEFAULT now(),
    "created_at" timestamptz DEFAULT now(),
    CONSTRAINT "album_file_album_fkey" FOREIGN KEY ("album") REFERENCES "public"."albums"("id") ON DELETE CASCADE,
    CONSTRAINT "album_file_added_by_fkey" FOREIGN KEY ("added_by") REFERENCES "public"."users"("id") ON DELETE SET NULL,
    CONSTRAINT "album_file_file_fkey" FOREIGN KEY ("file") REFERENCES "public"."files"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);
-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS user_album_id_seq;
-- Table Definition
CREATE TABLE IF NOT EXISTS "public"."user_album" (
  "id" int4 NOT NULL DEFAULT nextval('user_album_id_seq' :: regclass),
  "user" int4 NOT NULL,
  "album" int4 NOT NULL,
  "privilege" int2 NOT NULL DEFAULT 0,
  "updated_at" timestamptz DEFAULT now(),
  "created_at" timestamptz DEFAULT now(),
  CONSTRAINT "user_album_user_fkey" FOREIGN KEY ("user") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "user_album_album_fkey" FOREIGN KEY ("album") REFERENCES "public"."albums" ("id") ON DELETE CASCADE,
  PRIMARY KEY ("id")
);
-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS user_file_id_seq;
-- Table Definition
CREATE TABLE IF NOT EXISTS "public"."user_file" (
  "id" int4 NOT NULL DEFAULT nextval('user_file_id_seq' :: regclass),
  "user" int4 NOT NULL,
  "file" int4 NOT NULL,
  "privilege" int2 NOT NULL DEFAULT 0,
  "updated_at" timestamptz DEFAULT now(),
  "created_at" timestamptz DEFAULT now(),
  CONSTRAINT "user_file_user_fkey" FOREIGN KEY ("user") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "user_file_file_fkey" FOREIGN KEY ("file") REFERENCES "public"."files" ("id") ON DELETE CASCADE,
  PRIMARY KEY ("id")
);
