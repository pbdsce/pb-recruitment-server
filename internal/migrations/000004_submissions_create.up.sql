CREATE TYPE submission_type AS ENUM ('MCQ', 'Code');

CREATE TYPE submission_status AS ENUM (
    'pending',
    'failed_to_process',
    'accepted',
    'tle', 
    'mle',  
    'rte',  
    'wrong_answer'
);

CREATE TABLE submissions (
    id TEXT PRIMARY KEY, -- ULID
    user_id TEXT NOT NULL,       
    contest_id TEXT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    problem_id TEXT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    type submission_type NOT NULL,
    language TEXT,                
    code TEXT,                             
    choices INT[],                         
    status submission_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
