-- Remove existing duplicate drug plans for the same personal drug on the same Bangkok day.
-- Keep the earliest row and delete the rest.
WITH ranked AS (
  SELECT
    ctid,
    ROW_NUMBER() OVER (
      PARTITION BY pd_id, ((created_at AT TIME ZONE 'Asia/Bangkok')::date)
      ORDER BY created_at ASC, id ASC
    ) AS rn
  FROM drug_plans
)
DELETE FROM drug_plans dp
USING ranked r
WHERE dp.ctid = r.ctid
  AND r.rn > 1;

-- Enforce one plan per personal drug per Bangkok day.
CREATE UNIQUE INDEX IF NOT EXISTS idx_drug_plans_pd_id_bkk_day_unique
ON drug_plans (pd_id, ((created_at AT TIME ZONE 'Asia/Bangkok')::date));
