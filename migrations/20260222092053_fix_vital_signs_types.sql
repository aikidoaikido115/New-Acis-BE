-- Modify "vital_signs" table
ALTER TABLE "vital_signs" 
  ALTER COLUMN "temperature" TYPE numeric USING NULLIF(temperature, '')::numeric,
  ALTER COLUMN "heart_rate" TYPE smallint USING NULLIF(heart_rate, '')::smallint,
  ALTER COLUMN "breathing_rate" TYPE smallint USING NULLIF(breathing_rate, '')::smallint,
  ALTER COLUMN "blood_pressure_systolic" TYPE smallint USING NULLIF(blood_pressure_systolic, '')::smallint,
  ALTER COLUMN "blood_pressure_diastolic" TYPE smallint USING NULLIF(blood_pressure_diastolic, '')::smallint,
  ALTER COLUMN "oxygen_saturation" TYPE smallint USING NULLIF(oxygen_saturation, '')::smallint;
