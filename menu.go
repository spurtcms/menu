package menu

import (
	"regexp"
	"strings"
	"time"
)

var menumodel MenuModel

// MenuListing
func (menu *Menu) MenuList(limit int, offset int, filter Filter, tenantid string) (Menulist []TblMenus, count int64, err error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return []TblMenus{}, 0, AuthError
	}

	menumodel.DataAccess = menu.DataAccess
	menumodel.Userid = menu.UserId

	_, totalcount, _ := menumodel.MenuList(0, 0, filter, menu.DB, tenantid)

	menuparentlist, _, cerr := menumodel.MenuList(limit, offset, filter, menu.DB, tenantid)

	if cerr != nil {

		return []TblMenus{}, 0, cerr
	}

	return menuparentlist, totalcount, nil

}

// Create Menu
func (menu *Menu) CreateMenus(req MenuCreate) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}

	if req.MenuName == "" {

		return ErrorMenuName
	}

	var (
		menus     TblMenus
		menusSlug string
	)

	menus.Name = req.MenuName

	menusSlug = strings.ToLower(strings.ReplaceAll(req.MenuName, " ", "-"))

	menusSlug = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(menusSlug, "-")

	menusSlug = regexp.MustCompile(`-+`).ReplaceAllString(menusSlug, "-")

	menusSlug = strings.Trim(menusSlug, "-")

	menus.Slug = menusSlug

	menus.Description = req.Description

	menus.CreatedBy = req.CreatedBy

	menus.ParentId = 0

	menus.TenantId = req.TenantId

	menus.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := menumodel.CreateMenus(&menus, menu.DB)

	if err != nil {

		return err
	}

	return nil

}
