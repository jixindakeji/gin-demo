package e

const (
	SUCCESS       = 20000
	ERROR         = 500
	InvalidParams = 400

	TokenInvalid = 10000
	TokenExpired = 10001
	TokenError   = 10002
	TokenFailed  = 10003
	TokenMissed  = 10004

	UserAddFailed      = 20001
	UserGetFailed      = 20002
	UserGetCountFailed = 20003
	UserDeleteFailed   = 20004
	UserUpdateFailed   = 20005
	UserDisabled       = 20006

	MenuAddFailed    = 30001
	MenuDeleteFailed = 30002
	MenuUpdateFailed = 30003
	MenuGetFailed    = 30004

	RoleAddFailed    = 40001
	RoleGetFailed    = 40002
	RoleUpdateFailed = 40003
	RoleDeleteFailed = 40004

	RoleMenuUpdateFailed = 50001
	RoleMenuGetFailed    = 50002

	UserRoleGetFailed    = 60001
	UserRoleUpdateFailed = 60002
	UserMenuGetFailed    = 60003

	AuditsGetFailed = 70001
)
