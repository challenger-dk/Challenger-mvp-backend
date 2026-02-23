-- Create GiST spatial index on facilities.coordinates for efficient bounding box queries
CREATE INDEX IF NOT EXISTS "idx_facilities_coordinates" ON "facilities" USING GIST ("coordinates");
