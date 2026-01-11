-- Add distance column to challenges table (nullable, numeric)
ALTER TABLE "challenges" ADD COLUMN "distance" DOUBLE PRECISION NULL;
