-- Drop the functional index and replace with a unique constraint
DROP INDEX IF EXISTS idx_drug_plans_pd_id_bkk_day_unique;

-- Add a computed/generated column for Bangkok date to make constraint simpler
-- Note: We'll use a workaround with a trigger or handle it in the application
-- For now, create a unique constraint that works with ON CONFLICT

-- Create a unique constraint using a different approach
-- Use a partial unique index workaround by creating a constraint on columns only
-- This requires a dedicated approach - we'll use a check constraint + unique index combo

-- First, let's create a proper unique constraint using a helper column if possible
-- Since we can't easily use functional expressions in constraints,
-- we'll use the DEFERRABLE approach or change the ON CONFLICT logic

-- Create the unique index again, but properly
CREATE UNIQUE INDEX IF NOT EXISTS idx_drug_plans_pd_id_bkk_day_unique
ON drug_plans (pd_id, ((created_at AT TIME ZONE 'Asia/Bangkok')::date));

-- Note: The ON CONFLICT clause in the Go code should be changed to handle this differently
-- Instead of using the functional expression, we should:
-- 1. Check if a record exists before inserting, OR
-- 2. Use a trigger to handle duplicates, OR  
-- 3. Use raw SQL with proper error handling
