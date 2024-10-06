CREATE TABLE IF NOT EXISTS refresh_tokens (
    token text PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at timestamp(0) with time zone NOT NULL,
    revoked_at timestamp(0) with time zone DEFAULT NULL
)