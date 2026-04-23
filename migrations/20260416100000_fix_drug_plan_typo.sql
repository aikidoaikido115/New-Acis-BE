-- Rename typo columns only when needed to avoid failing on already-fixed databases.
DO $$
BEGIN
	IF EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND table_name = 'drug_plans'
		  AND column_name = 'is_ommitted'
	)
	AND NOT EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND table_name = 'drug_plans'
		  AND column_name = 'is_omitted'
	) THEN
		ALTER TABLE "drug_plans" RENAME COLUMN "is_ommitted" TO "is_omitted";
	END IF;

	IF EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND table_name = 'drug_plans'
		  AND column_name = 'ommitted_reason'
	)
	AND NOT EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND table_name = 'drug_plans'
		  AND column_name = 'omitted_reason'
	) THEN
		ALTER TABLE "drug_plans" RENAME COLUMN "ommitted_reason" TO "omitted_reason";
	END IF;
END
$$;
