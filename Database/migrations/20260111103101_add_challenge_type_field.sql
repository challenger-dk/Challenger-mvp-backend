-- Add type column to challenges table
ALTER TABLE "challenges" ADD COLUMN "type" character varying(20) NOT NULL DEFAULT 'open-for-all';

-- Add check constraint for type values
ALTER TABLE "challenges" ADD CONSTRAINT "chk_challenges_type" CHECK ((type)::text = ANY ((ARRAY['open-for-all'::character varying, 'team-vs-team'::character varying, 'run-cycling'::character varying])::text[]));
