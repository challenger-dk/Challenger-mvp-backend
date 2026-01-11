-- Add participants column to challenges table (nullable, integer)
ALTER TABLE "challenges" ADD COLUMN "participants" bigint NULL;
