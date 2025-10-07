CREATE MATERIALIZED VIEW ranking_mv AS
SELECT
    contest_id,
    user_id,
    score,
    hidden,
    disqualified,
    shortlisted,
    RANK() OVER (PARTITION BY contest_id ORDER BY score DESC) AS rank
FROM rankings;
