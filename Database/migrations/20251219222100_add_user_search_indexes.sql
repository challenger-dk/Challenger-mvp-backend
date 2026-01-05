-- atlas:txmode none

-- This migration adds trigram indexes to the users table to improve search performance

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_fullname_trgm
ON users USING gin ((first_name || ' ' || last_name) gin_trgm_ops);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email_trgm
ON users USING gin (email gin_trgm_ops);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_first_name_trgm
ON users USING gin (first_name gin_trgm_ops);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_last_name_trgm
ON users USING gin (last_name gin_trgm_ops);
