package entities

import "time"

type WarehouseTransaction struct {
	ID               string     `json:"id" gorm:"primaryKey"`
	Code             string     `json:"code" gorm:"not null;uniqueIndex"`
	Type             string     `json:"type" gorm:"not null"`
	ItemID           *string    `json:"-"`
	ItemCode         string     `json:"itemCode" gorm:"not null"`
	ItemName         string     `json:"itemName" gorm:"not null"`
	Quantity         int        `json:"quantity" gorm:"not null"`
	OperatorUserID   string     `json:"-" gorm:"not null"`
	Operator         string     `json:"operator" gorm:"not null"`
	ApprovalStatus   string     `json:"approvalStatus" gorm:"not null;default:'รออนุมัติ'"`
	ApprovedByUserID *string    `json:"-"`
	ApprovedBy       *string    `json:"approvedBy"`
	ApprovedAt       *time.Time `json:"-"`
	RejectedByUserID *string    `json:"-"`
	RejectedBy       *string    `json:"rejectedBy"`
	RejectedAt       *time.Time `json:"-"`
	RejectionReason  *string    `json:"rejectionReason"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`

	Item         WarehouseItem `json:"-" gorm:"foreignKey:ItemID;references:ID"`
	OperatorUser User          `json:"-" gorm:"foreignKey:OperatorUserID;references:ID"`
}
