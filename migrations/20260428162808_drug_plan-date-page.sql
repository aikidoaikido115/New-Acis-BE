-- Drop index "idx_laboratory_values_measurement_date" from table: "laboratory_values"
DROP INDEX "idx_laboratory_values_measurement_date";
-- Drop index "idx_laboratory_values_resident_measurement_date" from table: "laboratory_values"
DROP INDEX "idx_laboratory_values_resident_measurement_date";
-- Drop index "idx_vital_signs_measurement_date" from table: "vital_signs"
DROP INDEX "idx_vital_signs_measurement_date";
-- Drop index "idx_vital_signs_resident_measurement_date" from table: "vital_signs"
DROP INDEX "idx_vital_signs_resident_measurement_date";
-- Modify "warehouse_items" table
ALTER TABLE "warehouse_items" DROP CONSTRAINT "chk_warehouse_items_category", DROP CONSTRAINT "chk_warehouse_items_minimum_quantity_non_negative", DROP CONSTRAINT "chk_warehouse_items_quantity_non_negative", DROP CONSTRAINT "uni_warehouse_items_code", ALTER COLUMN "quantity" TYPE bigint, ALTER COLUMN "minimum_quantity" TYPE bigint, ALTER COLUMN "created_at" DROP NOT NULL, ALTER COLUMN "updated_at" DROP NOT NULL;
-- Create index "idx_warehouse_items_code" to table: "warehouse_items"
CREATE UNIQUE INDEX "idx_warehouse_items_code" ON "warehouse_items" ("code");
-- Drop index "idx_warehouse_transactions_approval_status" from table: "warehouse_transactions"
DROP INDEX "idx_warehouse_transactions_approval_status";
-- Drop index "idx_warehouse_transactions_created_at" from table: "warehouse_transactions"
DROP INDEX "idx_warehouse_transactions_created_at";
-- Drop index "idx_warehouse_transactions_item_code" from table: "warehouse_transactions"
DROP INDEX "idx_warehouse_transactions_item_code";
-- Modify "warehouse_transactions" table
ALTER TABLE "warehouse_transactions" DROP CONSTRAINT "chk_warehouse_transactions_approval_status", DROP CONSTRAINT "chk_warehouse_transactions_quantity_positive", DROP CONSTRAINT "chk_warehouse_transactions_type", DROP CONSTRAINT "uni_warehouse_transactions_code", DROP CONSTRAINT "fk_warehouse_transactions_approved_by", DROP CONSTRAINT "fk_warehouse_transactions_item", DROP CONSTRAINT "fk_warehouse_transactions_operator", DROP CONSTRAINT "fk_warehouse_transactions_rejected_by", ALTER COLUMN "quantity" TYPE bigint, ALTER COLUMN "created_at" DROP NOT NULL, ALTER COLUMN "updated_at" DROP NOT NULL, ADD CONSTRAINT "fk_warehouse_transactions_item" FOREIGN KEY ("item_id") REFERENCES "warehouse_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_warehouse_transactions_operator_user" FOREIGN KEY ("operator_user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "idx_warehouse_transactions_code" to table: "warehouse_transactions"
CREATE UNIQUE INDEX "idx_warehouse_transactions_code" ON "warehouse_transactions" ("code");
-- Drop "doctor_orders" table
DROP TABLE "doctor_orders";
-- Drop "support_tickets" table
DROP TABLE "support_tickets";
