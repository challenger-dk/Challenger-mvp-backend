-- Correct renames (preserve existing user prefs)
ALTER TABLE "user_settings" RENAME COLUMN "notify_team_invite" TO "notify_team_invites";
ALTER TABLE "user_settings" RENAME COLUMN "notify_friend_req" TO "notify_friend_requests";
ALTER TABLE "user_settings" RENAME COLUMN "notify_challenge_invite" TO "notify_challenge_invites";
ALTER TABLE "user_settings" RENAME COLUMN "notify_challenge_update" TO "notify_challenge_updates";

-- Add NEW toggles (these didn't exist before)
ALTER TABLE "user_settings"
  ADD COLUMN "notify_team_membership" boolean NOT NULL DEFAULT true,
  ADD COLUMN "notify_friend_updates" boolean NOT NULL DEFAULT true,
  ADD COLUMN "notify_challenge_reminders" boolean NOT NULL DEFAULT true;
