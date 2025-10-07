CREATE TABLE contests (
    id TEXT PRIMARY KEY,  -- UUID
    name TEXT NOT NULL,
    registration_start_time TIMESTAMP NOT NULL,
    registration_end_time TIMESTAMP NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
);