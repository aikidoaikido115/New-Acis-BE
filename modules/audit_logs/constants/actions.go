package constants

const (
	AuditActionInsert = "CREATE"
	AuditActionRead   = "READ"
	AuditActionUpdate = "UPDATE"
	AuditActionDelete = "DELETE"

	// AuditOldNewValuePassword is used to mask password values in audit logs
	AuditOldNewValuePassword = "********"
)
