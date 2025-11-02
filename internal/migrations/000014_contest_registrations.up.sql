CREATE TABLE contest_registrations (
    contest_id TEXT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    registered_at BIGINT NOT NULL,
    PRIMARY KEY (contest_id, user_id)
);
