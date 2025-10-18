ALTER TABLE submissions
ADD COLUMN runtime BIGINT,           
ADD COLUMN memory BIGINT,            
ADD COLUMN test_case_results JSONB; 