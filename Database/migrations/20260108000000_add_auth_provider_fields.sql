-- Add auth provider fields to users table
-- Make password nullable for OAuth users
ALTER TABLE "users" ALTER COLUMN "password" DROP NOT NULL;
ALTER TABLE "users" ADD COLUMN "auth_provider" text NOT NULL DEFAULT '';

-- Create index for auth provider field
CREATE INDEX IF NOT EXISTS "idx_users_auth_provider" ON "users" ("auth_provider");

