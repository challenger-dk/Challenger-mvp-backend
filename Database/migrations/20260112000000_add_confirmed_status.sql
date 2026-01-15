-- Add 'pending' and 'confirmed' status to challenges status constraint
ALTER TABLE "challenges" DROP CONSTRAINT IF EXISTS "chk_challenges_status";
ALTER TABLE "challenges" ADD CONSTRAINT "chk_challenges_status" CHECK ((status)::text = ANY ((ARRAY['suggested'::character varying, 'open'::character varying, 'pending'::character varying, 'ready'::character varying, 'confirmed'::character varying, 'completed'::character varying, 'exceeded'::character varying])::text[]));
