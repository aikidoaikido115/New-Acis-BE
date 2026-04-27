-- Create "warehouse_items" table
CREATE TABLE "warehouse_items" (
  "id" text NOT NULL,
  "code" text NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "quantity" integer NOT NULL DEFAULT 0,
  "minimum_quantity" integer NOT NULL DEFAULT 0,
  "unit" text NOT NULL,
  "category" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_warehouse_items_code" UNIQUE ("code"),
  CONSTRAINT "chk_warehouse_items_quantity_non_negative" CHECK ("quantity" >= 0),
  CONSTRAINT "chk_warehouse_items_minimum_quantity_non_negative" CHECK ("minimum_quantity" >= 0),
  CONSTRAINT "chk_warehouse_items_category" CHECK ("category" IN ('MED', 'EQU', 'CON'))
);

-- Create "warehouse_transactions" table
CREATE TABLE "warehouse_transactions" (
  "id" text NOT NULL,
  "code" text NOT NULL,
  "type" text NOT NULL,
  "item_id" text NULL,
  "item_code" text NOT NULL,
  "item_name" text NOT NULL,
  "quantity" integer NOT NULL,
  "operator_user_id" text NOT NULL,
  "operator" text NOT NULL,
  "approval_status" text NOT NULL DEFAULT 'รออนุมัติ',
  "approved_by_user_id" text NULL,
  "approved_by" text NULL,
  "approved_at" timestamptz NULL,
  "rejected_by_user_id" text NULL,
  "rejected_by" text NULL,
  "rejected_at" timestamptz NULL,
  "rejection_reason" text NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_warehouse_transactions_code" UNIQUE ("code"),
  CONSTRAINT "chk_warehouse_transactions_quantity_positive" CHECK ("quantity" > 0),
  CONSTRAINT "chk_warehouse_transactions_type" CHECK ("type" IN ('เพิ่มสินค้าใหม่', 'เติมสินค้า', 'เบิกสินค้า', 'นำออก')),
  CONSTRAINT "chk_warehouse_transactions_approval_status" CHECK ("approval_status" IN ('รออนุมัติ', 'อนุมัติ', 'ไม่อนุมัติ')),
  CONSTRAINT "fk_warehouse_transactions_item" FOREIGN KEY ("item_id") REFERENCES "warehouse_items" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_warehouse_transactions_operator" FOREIGN KEY ("operator_user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_warehouse_transactions_approved_by" FOREIGN KEY ("approved_by_user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_warehouse_transactions_rejected_by" FOREIGN KEY ("rejected_by_user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);

CREATE INDEX "idx_warehouse_transactions_item_code" ON "warehouse_transactions" ("item_code");
CREATE INDEX "idx_warehouse_transactions_approval_status" ON "warehouse_transactions" ("approval_status");
CREATE INDEX "idx_warehouse_transactions_created_at" ON "warehouse_transactions" ("created_at");
