package enums

type AuditLogAction string

const (
	AuditLogActionCreate AuditLogAction = "create"
	AuditLogActionUpdate AuditLogAction = "update"
	AuditLogActionDelete AuditLogAction = "delete"
	AuditLogActionLogin  AuditLogAction = "login"
	AuditLogActionLogout AuditLogAction = "logout"
)
