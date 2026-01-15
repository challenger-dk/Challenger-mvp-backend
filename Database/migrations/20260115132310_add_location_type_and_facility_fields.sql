-- Add location type field with default value 'address'
ALTER TABLE "locations" ADD COLUMN IF NOT EXISTS "type" VARCHAR(20) NOT NULL DEFAULT 'address';
ALTER TABLE "locations" ADD CONSTRAINT "chk_locations_type" CHECK (type IN ('address', 'facility'));

-- Add facility-specific fields (all nullable)
ALTER TABLE "locations" ADD COLUMN IF NOT EXISTS "facility_id" text NULL;
ALTER TABLE "locations" ADD COLUMN IF NOT EXISTS "facility_name" text NULL;
ALTER TABLE "locations" ADD COLUMN IF NOT EXISTS "detailed_name" text NULL;
ALTER TABLE "locations" ADD COLUMN IF NOT EXISTS "email" text NULL;
ALTER TABLE "locations" ADD COLUMN IF NOT EXISTS "website" text NULL;
ALTER TABLE "locations" ADD COLUMN IF NOT EXISTS "facility_type" text NULL;
ALTER TABLE "locations" ADD COLUMN IF NOT EXISTS "indoor" boolean NULL;
ALTER TABLE "locations" ADD COLUMN IF NOT EXISTS "notes" text NULL;
