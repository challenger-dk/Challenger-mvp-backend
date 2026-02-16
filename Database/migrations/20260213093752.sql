-- Create "facilities" table
CREATE TABLE "facilities" (
  "id" bigserial NOT NULL,
  "external_id" text NOT NULL,
  "name" text NOT NULL,
  "detailed_name" text NOT NULL,
  "address" text NOT NULL,
  "website" text NULL,
  "email" text NULL,
  "facility_type" text NOT NULL,
  "indoor" boolean NULL DEFAULT false,
  "notes" text NULL,
  "coordinates" geography(Point,4326) NOT NULL,
  "postal_code" text NOT NULL,
  "city" text NOT NULL,
  "country" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_facilities_external_id" to table: "facilities"
CREATE UNIQUE INDEX "idx_facilities_external_id" ON "facilities" ("external_id");
-- Modify "challenges" table
ALTER TABLE "challenges" ADD COLUMN "facility_id" bigint NULL, ADD CONSTRAINT "fk_facilities_challenges" FOREIGN KEY ("facility_id") REFERENCES "facilities" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
