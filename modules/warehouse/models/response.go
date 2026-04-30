package models

type WarehouseTransactionResponse struct {
	ID              string  `json:"id"`
	Code            string  `json:"code"`
	Type            string  `json:"type"`
	ItemCode        string  `json:"itemCode"`
	ItemName        string  `json:"itemName"`
	Quantity        int     `json:"quantity"`
	Operator        string  `json:"operator"`
	Date            string  `json:"date"`
	ApprovalStatus  string  `json:"approvalStatus"`
	ApprovedBy      *string `json:"approvedBy,omitempty"`
	ApprovedAt      *string `json:"approvedAt,omitempty"`
	RejectedBy      *string `json:"rejectedBy,omitempty"`
	RejectedAt      *string `json:"rejectedAt,omitempty"`
	RejectionReason *string `json:"rejectionReason,omitempty"`
}
