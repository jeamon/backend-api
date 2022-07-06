BEGIN;

CREATE TABLE IF NOT EXISTS data.scan_infos
(
    id UUID DEFAULT uuid_generate_v1mc() PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    company_id TEXT NOT NULL,
    username TEXT NOT NULL DEFAULT '',
    client_id TEXT NOT NULL,
    repository_url TEXT NOT NULL,
    commit_id TEXT NOT NULL DEFAULT '',
    tag_id TEXT NOT NULL DEFAULT '',
    results TEXT [] NOT NULL DEFAULT ARRAY []::TEXT [],
    started_at bigserial NOT NULL,
    completed_at bigserial NOT NULL,
    sent_at bigserial NOT NULL,
    error TEXT NOT NULL DEFAULT '',
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb
);

COMMIT;