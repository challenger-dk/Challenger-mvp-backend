-- Modify "challenges" table
ALTER TABLE "challenges" ADD COLUMN "tags" jsonb NULL DEFAULT '[]';
