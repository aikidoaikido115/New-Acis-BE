-- Rename role "Super Admin" to "Super User"
UPDATE "roles"
SET "name" = 'Super User'
WHERE "name" = 'Super Admin';
