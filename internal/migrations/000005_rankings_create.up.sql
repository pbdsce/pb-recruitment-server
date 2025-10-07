CREATE TABLE rankings (
    contest_id TEXT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    score INT NOT NULL DEFAULT 0,
    hidden BOOLEAN NOT NULL DEFAULT FALSE,
    disqualified BOOLEAN NOT NULL DEFAULT FALSE,
    shortlisted BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (contest_id, user_id)
);
