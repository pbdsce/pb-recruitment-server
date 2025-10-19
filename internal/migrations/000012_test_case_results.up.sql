ALTER TABLE submissions
DROP COLUMN test_case_results;

CREATE TABLE test_case_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id TEXT NOT NULL,
    test_case_id TEXT NOT NULL, 
    status TEXT NOT NULL,
    runtime BIGINT NOT NULL,
    memory BIGINT NOT NULL,
    created_at BIGINT NOT NULL,

    CONSTRAINT fk_submission
        FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE
);