package e

//MsgFlags global msg flags
var MsgFlags = map[int]string{
	SUCCESS:       "OK",
	ERROR:         "fail",
	InvalidParams: "请求参数错误",

	TokenInvalid: "Token不合法",
	TokenExpired: "Token过期",
	TokenError:   "Token错误",
	TokenFailed:  "token失败",
	TokenMissed:  "token不存在",

	UserAddFailed:      "新增用户失败",
	UserGetFailed:      "获取用户失败",
	UserGetCountFailed: "获取用户数量失败",
	UserDeleteFailed:   "删除用户失败",
	UserUpdateFailed:   "更新用户失败",
	UserDisabled:       "用户被禁止登录",

	MenuAddFailed:    "新增菜单失败",
	MenuDeleteFailed: "删除菜单失败",
	MenuUpdateFailed: "修改菜单失败",
	MenuGetFailed:    "获取菜单失败",

	RoleAddFailed:    "添加角色失败",
	RoleGetFailed:    "获取角色失败",
	RoleUpdateFailed: "更新角色失败",
	RoleDeleteFailed: "删除角色失败",

	RoleMenuGetFailed:    "获取角色菜单失败",
	RoleMenuUpdateFailed: "更新角色菜单失败",

	UserRoleGetFailed:    "获取用户角色失败",
	UserRoleUpdateFailed: "更新用户角色失败",
	UserMenuGetFailed:    "获取用户菜单失败",

	AuditsGetFailed: "获取操作记录失败",
}

//GetMsg get global msg by code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
