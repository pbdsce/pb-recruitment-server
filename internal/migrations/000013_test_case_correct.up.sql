ALTER TABLE submissions
DROP COLUMN IF EXISTS test_case_results;

CREATE TYPE test_case_status AS ENUM (
    'pass', 
    'wrong_answer', 
    'tle', 
    'mle', 
    'rte'
);

CREATE TABLE IF NOT EXISTS test_case_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id TEXT NOT NULL,
    test_case_id TEXT NOT NULL, 
    status test_case_status NOT NULL,
    runtime BIGINT NOT NULL,
    memory BIGINT NOT NULL,
    created_at BIGINT NOT NULL,

    CONSTRAINT fk_submission
        FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE
);