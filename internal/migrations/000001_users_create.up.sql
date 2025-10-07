CREATE TABLE users (
    id TEXT PRIMARY KEY,  -- Firebase UID
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    usn TEXT NOT NULL UNIQUE,
    mobile_number TEXT,
    joining_year INT NOT NULL,
    department TEXT NOT NULL
);