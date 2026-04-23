package constants

const (
	CategoryMedical   = "MED"
	CategoryEquipment = "EQU"
	CategoryConsumer  = "CON"
)

const (
	TransactionTypeAddItem  = "เพิ่มสินค้าใหม่"
	TransactionTypeRestock  = "เติมสินค้า"
	TransactionTypeWithdraw = "เบิกสินค้า"
	TransactionTypeRemove   = "นำออก"
)

const (
	ApprovalStatusPending  = "รออนุมัติ"
	ApprovalStatusApproved = "อนุมัติ"
	ApprovalStatusRejected = "ไม่อนุมัติ"
)

const (
	AdjustModeRestock  = "restock"
	AdjustModeWithdraw = "withdraw"
)
