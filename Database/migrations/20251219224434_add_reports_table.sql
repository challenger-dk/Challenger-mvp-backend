-- Create "reports" table
CREATE TABLE "reports" (
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
