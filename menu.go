package menu

import (
	"regexp"
	"strings"
	"time"

	"github.com/spurtcms/menu/migration"
)

func MenuSetup(config Config) *Menu {

	migration.AutoMigration(config.DB, config.DataBaseType)

	return &Menu{
		DB:               config.DB,
		AuthEnable:       config.AuthEnable,
		PermissionEnable: config.PermissionEnable,
		Auth:             config.Auth,
		Permissions:      config.Permissions,
	}

}

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

	menus.SlugName = menusSlug

	menus.Description = req.Description

	menus.CreatedBy = req.CreatedBy

	menus.ParentId = 0

	menus.TenantId = req.TenantId

	menus.Status = req.Status

	menus.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	menus.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := menumodel.CreateMenus(&menus, menu.DB)

	if err != nil {

		return err
	}

	return nil

}

/*UpdateMenu*/
func (menu *Menu) UpdateMenu(req MenuCreate) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}

	if req.Id <= 0 || req.MenuName == "" {

		return ErrorMenuName
	}

	var (
		menudet     TblMenus
		menudetSlug string
	)

	menudet.Id = req.Id

	menudet.Name = req.MenuName

	menudet.Description = req.Description

	menudetSlug = strings.ToLower(strings.ReplaceAll(req.MenuName, " ", "-"))

	menudetSlug = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(menudetSlug, "-")

	menudetSlug = regexp.MustCompile(`-+`).ReplaceAllString(menudetSlug, "-")

	menudetSlug = strings.Trim(menudetSlug, "-")

	menudet.Status = req.Status

	menudet.SlugName = menudetSlug

	menudet.ModifiedBy = req.ModifiedBy

	menudet.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	menudet.TenantId = req.TenantId

	err := menumodel.UpdateMenu(&menudet, menu.DB)

	if err != nil {

		return err
	}

	return nil

}

/*Deletemenu*/
func (menu *Menu) DeleteMenu(menuid int, modifiedby int, tenantid string) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}
	GetData, _ := menumodel.GetMenuTree(menuid, menu.DB, tenantid)

	var individualid []int

	for _, GetParent := range GetData {

		indivi := GetParent.Id

		individualid = append(individualid, indivi)
	}

	var menudet TblMenus

	menudet.DeletedBy = modifiedby

	menudet.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	menudet.IsDeleted = 1

	err := menumodel.DeleteMenuById(&menudet, individualid, tenantid, menu.DB)

	if err != nil {

		return err
	}

	return nil

}

// Check Menuname is already exists
func (menu *Menu) CheckMenuName(id int, name string, tenantid string) (bool, error) {

	var menudet TblMenus

	err := menumodel.CheckMenuName(menudet, id, name, menu.DB, tenantid)

	if err != nil {

		return false, err

	}

	return true, nil
}
