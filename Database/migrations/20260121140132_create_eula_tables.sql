-- Modify "challenges" table
ALTER TABLE "challenges" ALTER COLUMN "distance" TYPE numeric;
-- Modify "conversations" table
ALTER TABLE "conversations" DROP CONSTRAINT "chk_conversations_type", ADD CONSTRAINT "chk_conversations_type" CHECK ((type)::text = ANY ((ARRAY['direct'::character varying, 'group'::character varying, 'team'::character varying, 'challenge'::character varying])::text[]));
-- Modify "user_settings" table
ALTER TABLE "user_settings" ALTER COLUMN "notify_team_membership" DROP NOT NULL, ALTER COLUMN "notify_friend_updates" DROP NOT NULL, ALTER COLUMN "notify_challenge_reminders" DROP NOT NULL;
-- Drop index "idx_users_auth_provider" from table: "users"
DROP INDEX "idx_users_auth_provider";
-- Modify "users" table
ALTER TABLE "users" ALTER COLUMN "auth_provider" DROP NOT NULL;
-- Create "eula_acceptances" table
CREATE TABLE "eula_acceptances" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "eula_version_id" bigint NOT NULL,
  "accepted_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_user_eula" to table: "eula_acceptances"
CREATE UNIQUE INDEX "idx_user_eula" ON "eula_acceptances" ("user_id", "eula_version_id");
-- Create "eula_versions" table
CREATE TABLE "eula_versions" (
  "id" bigserial NOT NULL,
  "version" text NOT NULL,
  "locale" text NOT NULL,
  "content" text NOT NULL,
  "content_hash" character(64) NOT NULL,
  "is_active" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_eula_version_locale" to table: "eula_versions"
CREATE UNIQUE INDEX "idx_eula_version_locale" ON "eula_versions" ("version", "locale");
