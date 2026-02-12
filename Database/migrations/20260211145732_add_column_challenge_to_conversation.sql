-- Modify "conversations" table
ALTER TABLE "conversations" ADD COLUMN "challenge_id" bigint NULL, ADD CONSTRAINT "fk_conversations_challenge" FOREIGN KEY ("challenge_id") REFERENCES "challenges" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "idx_challenge_conversation" to table: "conversations"
CREATE UNIQUE INDEX "idx_challenge_conversation" ON "conversations" ("challenge_id");
