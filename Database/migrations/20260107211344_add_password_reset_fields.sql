-- Add password reset fields to users table
ALTER TABLE "users" ADD COLUMN "password_reset_code" text NULL;
ALTER TABLE "users" ADD COLUMN "password_reset_code_expires_at" timestamptz NULL;

-- Create indexes for password reset fields
CREATE INDEX IF NOT EXISTS "idx_users_password_reset_code" ON "users" ("password_reset_code");
CREATE INDEX IF NOT EXISTS "idx_users_password_reset_code_expires_at" ON "users" ("password_reset_code_expires_at");

