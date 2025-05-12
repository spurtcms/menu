package menu

import "errors"

var (
	ErrorAuth         = errors.New("auth enabled not initialised")
	ErrorPermission   = errors.New("permissions enabled not initialised")
	ErrorMenuName = errors.New("given some values is empty")
)

func AuthandPermission(menu *Menu) error {

	//check auth enable if enabled, use auth pkg otherwise it will return error
	if menu.AuthEnable && !menu.Auth.AuthFlg {

		return ErrorAuth
	}
	//check permission enable if enabled, use team-role pkg otherwise it will return error
	if menu.PermissionEnable && !menu.Auth.PermissionFlg {

		return ErrorPermission

	}

	return nil
}
