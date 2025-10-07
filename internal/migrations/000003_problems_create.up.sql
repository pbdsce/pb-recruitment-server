CREATE TYPE problem_type AS ENUM ('code', 'mcq');

CREATE TABLE problems (
    id TEXT PRIMARY KEY,
    contest_id TEXT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    score INT NOT NULL,
    type problem_type NOT NULL,
    answer INT[]
);
