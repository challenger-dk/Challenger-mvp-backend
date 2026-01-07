-- Create required extensions
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create "locations" table
CREATE TABLE IF NOT EXISTS "locations" (
  "id" bigserial NOT NULL,
  "address" text NOT NULL,
  "coordinates" geography(Point,4326) NOT NULL,
  "postal_code" text NOT NULL,
  "city" text NOT NULL,
  "country" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_locations_coordinates" to table: "locations"
CREATE UNIQUE INDEX IF NOT EXISTS "idx_locations_coordinates" ON "locations" ("coordinates");
-- Create "users" table
CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial NOT NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "first_name" text NOT NULL,
  "last_name" text NULL,
  "profile_picture" text NULL,
  "bio" text NULL,
  "birth_date" timestamptz NULL,
  "city" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX IF NOT EXISTS "idx_users_deleted_at" ON "users" ("deleted_at");
-- Create "challenges" table
CREATE TABLE IF NOT EXISTS "challenges" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "sport" text NULL,
  "location_id" bigint NOT NULL,
  "creator_id" bigint NOT NULL,
  "is_indoor" boolean NULL DEFAULT false,
  "is_public" boolean NULL DEFAULT false,
  "is_completed" boolean NULL DEFAULT false,
  "status" character varying(20) NOT NULL DEFAULT 'open',
  "play_for" text NULL,
  "has_cost" boolean NULL DEFAULT false,
  "comment" text NULL,
  "team_size" bigint NULL,
  "date" timestamptz NOT NULL,
  "start_time" timestamptz NOT NULL,
  "end_time" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_challenges_location" FOREIGN KEY ("location_id") REFERENCES "locations" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_users_created_challenges" FOREIGN KEY ("creator_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_challenges_status" CHECK ((status)::text = ANY ((ARRAY['suggested'::character varying, 'open'::character varying, 'ready'::character varying, 'completed'::character varying, 'exceeded'::character varying])::text[]))
);
-- Create index "idx_challenges_deleted_at" to table: "challenges"
CREATE INDEX IF NOT EXISTS "idx_challenges_deleted_at" ON "challenges" ("deleted_at");
-- Create "teams" table
CREATE TABLE IF NOT EXISTS "teams" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "creator_id" bigint NOT NULL,
  "location_id" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_teams_creator" FOREIGN KEY ("creator_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_teams_location" FOREIGN KEY ("location_id") REFERENCES "locations" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_teams_deleted_at" to table: "teams"
CREATE INDEX IF NOT EXISTS "idx_teams_deleted_at" ON "teams" ("deleted_at");
-- Create index "idx_teams_location_id" to table: "teams"
CREATE INDEX IF NOT EXISTS "idx_teams_location_id" ON "teams" ("location_id");
-- Create "challenge_teams" table
CREATE TABLE IF NOT EXISTS "challenge_teams" (
  "challenge_id" bigint NOT NULL,
  "team_id" bigint NOT NULL,
  PRIMARY KEY ("challenge_id", "team_id"),
  CONSTRAINT "fk_challenge_teams_challenge" FOREIGN KEY ("challenge_id") REFERENCES "challenges" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_challenge_teams_team" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "conversations" table
CREATE TABLE IF NOT EXISTS "conversations" (
  "id" bigserial NOT NULL,
  "type" character varying(20) NOT NULL,
  "title" character varying(255) NULL,
  "team_id" bigint NULL,
  "direct_key" character varying(100) NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_conversations_team" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_conversations_type" CHECK ((type)::text = ANY ((ARRAY['direct'::character varying, 'group'::character varying, 'team'::character varying])::text[]))
);
-- Create index "idx_direct_key" to table: "conversations"
CREATE UNIQUE INDEX IF NOT EXISTS "idx_direct_key" ON "conversations" ("direct_key");
-- Create index "idx_team_conversation" to table: "conversations"
CREATE UNIQUE INDEX IF NOT EXISTS "idx_team_conversation" ON "conversations" ("team_id");
-- Create "conversation_participants" table
CREATE TABLE IF NOT EXISTS "conversation_participants" (
  "conversation_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "joined_at" timestamptz NULL,
  "last_read_at" timestamptz NULL,
  "left_at" timestamptz NULL,
  PRIMARY KEY ("conversation_id", "user_id"),
  CONSTRAINT "fk_conversation_participants_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_conversations_participants" FOREIGN KEY ("conversation_id") REFERENCES "conversations" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "emergency_infos" table
CREATE TABLE IF NOT EXISTS "emergency_infos" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "name" text NOT NULL,
  "phone_number" text NOT NULL,
  "relationship" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_emergency_contacts" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create "invitations" table
CREATE TABLE IF NOT EXISTS "invitations" (
  "id" bigserial NOT NULL,
  "inviter_id" bigint NOT NULL,
  "invitee_id" bigint NOT NULL,
  "note" text NULL,
  "resource_type" character varying(20) NOT NULL,
  "resource_id" bigint NOT NULL,
  "status" character varying(20) NOT NULL DEFAULT 'pending',
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_invitations_invitee" FOREIGN KEY ("invitee_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_invitations_inviter" FOREIGN KEY ("inviter_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_invitations_resource_type" CHECK ((resource_type)::text = ANY ((ARRAY['team'::character varying, 'friend'::character varying, 'challenge'::character varying])::text[])),
  CONSTRAINT "chk_invitations_status" CHECK ((status)::text = ANY ((ARRAY['pending'::character varying, 'accepted'::character varying, 'declined'::character varying])::text[]))
);
-- Create index "idx_unique_invitation" to table: "invitations"
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_invitation" ON "invitations" ("inviter_id", "invitee_id", "resource_type", "resource_id");
-- Create "messages" table
CREATE TABLE IF NOT EXISTS "messages" (
  "id" bigserial NOT NULL,
  "conversation_id" bigint NULL,
  "sender_id" bigint NOT NULL,
  "team_id" bigint NULL,
  "recipient_id" bigint NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_conversations_messages" FOREIGN KEY ("conversation_id") REFERENCES "conversations" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_messages_recipient" FOREIGN KEY ("recipient_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_messages_sender" FOREIGN KEY ("sender_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_conversation_created" to table: "messages"
CREATE INDEX IF NOT EXISTS "idx_conversation_created" ON "messages" ("conversation_id", "created_at");
-- Create index "idx_messages_recipient_id" to table: "messages"
CREATE INDEX IF NOT EXISTS "idx_messages_recipient_id" ON "messages" ("recipient_id");
-- Create index "idx_messages_team_id" to table: "messages"
CREATE INDEX IF NOT EXISTS "idx_messages_team_id" ON "messages" ("team_id");
-- Create "notifications" table
CREATE TABLE IF NOT EXISTS "notifications" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "actor_id" bigint NULL,
  "type" character varying(50) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "resource_id" bigint NULL,
  "resource_type" text NULL,
  "invitation_id" bigint NULL,
  "is_read" boolean NULL DEFAULT false,
  "is_relevant" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_notifications_actor" FOREIGN KEY ("actor_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_notifications_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_notifications_actor_id" to table: "notifications"
CREATE INDEX IF NOT EXISTS "idx_notifications_actor_id" ON "notifications" ("actor_id");
-- Create index "idx_notifications_user_id" to table: "notifications"
CREATE INDEX IF NOT EXISTS "idx_notifications_user_id" ON "notifications" ("user_id");
-- Create "reports" table
CREATE TABLE IF NOT EXISTS "reports" (
  "id" bigserial NOT NULL,
  "reporter_id" bigint NOT NULL,
  "target_id" bigint NOT NULL,
  "target_type" text NOT NULL,
  "reason" text NOT NULL,
  "comment" text NULL,
  "status" text NULL DEFAULT 'PENDING',
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_reports_reporter" FOREIGN KEY ("reporter_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "sports" table
CREATE TABLE IF NOT EXISTS "sports" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_sports_name" UNIQUE ("name")
);
-- Create "team_sports" table
CREATE TABLE IF NOT EXISTS "team_sports" (
  "team_id" bigint NOT NULL,
  "sport_id" bigint NOT NULL,
  PRIMARY KEY ("team_id", "sport_id"),
  CONSTRAINT "fk_team_sports_sport" FOREIGN KEY ("sport_id") REFERENCES "sports" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_team_sports_team" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "user_blocked_users" table
CREATE TABLE IF NOT EXISTS "user_blocked_users" (
  "user_id" bigint NOT NULL,
  "blocked_user_id" bigint NOT NULL,
  PRIMARY KEY ("user_id", "blocked_user_id"),
  CONSTRAINT "fk_user_blocked_users_blocked_users" FOREIGN KEY ("blocked_user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_user_blocked_users_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "user_challenges" table
CREATE TABLE IF NOT EXISTS "user_challenges" (
  "user_id" bigint NOT NULL,
  "challenge_id" bigint NOT NULL,
  PRIMARY KEY ("user_id", "challenge_id"),
  CONSTRAINT "fk_user_challenges_challenge" FOREIGN KEY ("challenge_id") REFERENCES "challenges" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_user_challenges_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "user_favorite_sports" table
CREATE TABLE IF NOT EXISTS "user_favorite_sports" (
  "user_id" bigint NOT NULL,
  "sport_id" bigint NOT NULL,
  PRIMARY KEY ("user_id", "sport_id"),
  CONSTRAINT "fk_user_favorite_sports_sport" FOREIGN KEY ("sport_id") REFERENCES "sports" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_user_favorite_sports_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "user_friends" table
CREATE TABLE IF NOT EXISTS "user_friends" (
  "user_id" bigint NOT NULL,
  "friend_id" bigint NOT NULL,
  PRIMARY KEY ("user_id", "friend_id"),
  CONSTRAINT "fk_user_friends_friends" FOREIGN KEY ("friend_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_user_friends_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "user_settings" table
CREATE TABLE IF NOT EXISTS "user_settings" (
  "user_id" bigserial NOT NULL,
  "notify_team_invite" boolean NULL DEFAULT true,
  "notify_friend_req" boolean NULL DEFAULT true,
  "notify_challenge_invite" boolean NULL DEFAULT true,
  "notify_challenge_update" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("user_id"),
  CONSTRAINT "fk_users_settings" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create "user_teams" table
CREATE TABLE IF NOT EXISTS "user_teams" (
  "team_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("team_id", "user_id"),
  CONSTRAINT "fk_user_teams_team" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_user_teams_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
