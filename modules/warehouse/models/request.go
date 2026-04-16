package models

type ListWarehouseItemsQuery struct {
	Search   string `query:"search" json:"search"`
	Category string `query:"category" json:"category"`
}

type CreateWarehouseItemRequest struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Quantity        int    `json:"quantity"`
	MinimumQuantity int    `json:"minimumQuantity"`
	Unit            string `json:"unit"`
	Category        string `json:"category"`
}

type UpdateWarehouseItemRequest struct {
	Code            *string `json:"code"`
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	Quantity        *int    `json:"quantity"`
	MinimumQuantity *int    `json:"minimumQuantity"`
	Unit            *string `json:"unit"`
	Category        *string `json:"category"`
}

type AdjustWarehouseItemRequest struct {
	Mode     string `json:"mode"`
	Quantity int    `json:"quantity"`
}

type ListWarehouseTransactionsQuery struct {
	StartDate  string `query:"startDate" json:"startDate"`
	EndDate    string `query:"endDate" json:"endDate"`
	SearchItem string `query:"searchItem" json:"searchItem"`
	SearchUser string `query:"searchUser" json:"searchUser"`
	Status     string `query:"status" json:"status"`
	Type       string `query:"type" json:"type"`
}

type ApproveTransactionsRequest struct {
	TransactionIDs []string `json:"transactionIds"`
}

type RejectTransactionsRequest struct {
	TransactionIDs []string `json:"transactionIds"`
	Reason         string   `json:"reason"`
}
