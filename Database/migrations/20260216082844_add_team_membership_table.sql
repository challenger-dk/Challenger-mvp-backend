-- Create "team_members" table
CREATE TABLE "team_members" (
  "team_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "role" character varying(20) NOT NULL DEFAULT 'member',
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("team_id", "user_id"),
  CONSTRAINT "fk_team_members_team" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_team_members_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);

-- Create index on team_id
CREATE INDEX "idx_team_members_team_id" ON "team_members" ("team_id");

-- Create index on user_id
CREATE INDEX "idx_team_members_user_id" ON "team_members" ("user_id");

-- Create index on role
CREATE INDEX "idx_team_members_role" ON "team_members" ("role");

-- Migrate data from user_teams to team_members
-- Set creators as owners, others as members
INSERT INTO "team_members" ("team_id", "user_id", "role", "created_at")
SELECT
  ut.team_id,
  ut.user_id,
  CASE
    WHEN t.creator_id = ut.user_id THEN 'owner'
    ELSE 'member'
  END as role,
  CURRENT_TIMESTAMP
FROM "user_teams" ut
JOIN "teams" t ON ut.team_id = t.id;

-- Drop "user_teams" table
DROP TABLE "user_teams";
