package models

type CreateSupportTicketRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type ListSupportTicketsQuery struct {
	Search          string `query:"search" json:"search"`
	Status          string `query:"status" json:"status"`
	ReporterRole    string `query:"reporterRole" json:"reporterRole"`
	CreatedByUserID string `query:"createdByUserId" json:"createdByUserId"`
}

type UpdateSupportTicketStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
