CREATE TABLE contests (
    id TEXT PRIMARY KEY,  -- UUID
    name TEXT NOT NULL,
    registration_start_time BIGINT NOT NULL,
    registration_end_time BIGINT NOT NULL,
    start_time BIGINT NOT NULL,
    end_time BIGINT NOT NULL
);