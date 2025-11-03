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
func (menu *Menu) MenuList(limit int, offset int, filter Filter, tenantid string, websiteid int) (Menulist []TblMenus, count int64, err error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return []TblMenus{}, 0, AuthError
	}

	menumodel.DataAccess = menu.DataAccess
	menumodel.Userid = menu.UserId

	_, totalcount, _ := menumodel.MenuList(0, 0, filter, menu.DB, tenantid, websiteid)

	menuparentlist, _, cerr := menumodel.MenuList(limit, offset, filter, menu.DB, tenantid, websiteid)

	if cerr != nil {

		return []TblMenus{}, 0, cerr
	}

	return menuparentlist, totalcount, nil

}

// Create Menu
func (menu *Menu) CreateMenus(req MenuCreate) (TblMenus, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblMenus{}, AuthError
	}

	if req.MenuName == "" {

		return TblMenus{}, ErrorMenuName
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

	menus.ParentId = req.ParentId

	menus.TenantId = req.TenantId

	menus.Status = req.Status

	menus.UrlPath = req.UrlPath

	menus.Type = req.Type

	menus.TypeId = req.TypeId

	menus.WebsiteId = req.WebsiteId

	menus.ListingsIds = req.ListingsIds

	menus.CategoryIds = req.CategoryIds

	menus.ImageName = req.ImageName

	menus.ImagePath = req.ImagePath

	menus.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	menus.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	menn, err := menumodel.CreateMenus(&menus, menu.DB)

	if err != nil {

		return TblMenus{}, err
	}

	return menn, nil

}

/*UpdateMenu*/
func (menu *Menu) UpdateMenu(req MenuCreate) (TblMenus, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblMenus{}, AuthError
	}

	if req.Id <= 0 || req.MenuName == "" {

		return TblMenus{}, ErrorMenuName
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

	menudet.ParentId = req.ParentId

	menudet.UrlPath = req.UrlPath

	menudet.Type = req.Type

	menudet.TypeId = req.TypeId

	menudet.WebsiteId = req.WebsiteId

	menudet.ImageName = req.ImageName

	menudet.ImagePath = req.ImagePath

	menudet.CategoryIds = req.CategoryIds

	menudet.ListingsIds = req.ListingsIds

	updatemenu, err := menumodel.UpdateMenu(&menudet, menu.DB)

	if err != nil {

		return TblMenus{}, err
	}

	return updatemenu, nil

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
func (menu *Menu) CheckMenuName(id int, name string, websiteid int, tenantid string) (bool, error) {

	var menudet TblMenus

	err := menumodel.CheckMenuName(menudet, id, name, websiteid, menu.DB, tenantid)

	if err != nil {

		return false, err

	}

	return true, nil
}

//MenuStatuc Update Function//

func (menu *Menu) MenuStatusChange(menuid int, status int, userid int, tenantid string) (bool, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return false, AuthError
	}

	var menudet TblMenus

	menudet.ModifiedBy = userid

	menudet.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	menudet.Status = status

	menudet.TenantId = tenantid

	menudet.Id = menuid

	err := menumodel.MenuStatusChange(menudet, menu.DB)

	if err != nil {

		return false, err

	}
	return true, nil
}

func (menu *Menu) GetMenusByParentid(parentid int, tenantid string) ([]TblMenus, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return []TblMenus{}, AuthError
	}
	GetData, _ := menumodel.GetMenuTree(parentid, menu.DB, tenantid)

	return GetData, nil
}

func (menu *Menu) DeleteMenuItem(menuid int, userid int, tenantid string) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}
	var menudet TblMenus

	menudet.DeletedBy = userid

	menudet.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	menudet.IsDeleted = 1

	menudet.Id = menuid

	menudet.TenantId = tenantid

	err := menumodel.DeleteMenuItemById(menudet, menu.DB)

	if err != nil {

		return err
	}

	return nil
}

func (menu *Menu) GetmenyById(menuid int, tenantid string) (TblMenus, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblMenus{}, AuthError
	}

	GetData, _ := menumodel.GetMenuById(menuid, menu.DB, tenantid)

	return GetData, nil
}

func (menu *Menu) GetMenuBySlug(slug, tenantid string) (TblMenus, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblMenus{}, AuthError
	}

	GetData, _ := menumodel.GetMenuBySlug(slug, menu.DB, tenantid)

	return GetData, nil
}

func (menu *Menu) GetMenuBySlugName(slug string, websiteid int, tenantid string) (TblMenus, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblMenus{}, AuthError
	}

	GetData, _ := menumodel.GetMenuBySlugName(slug, websiteid, menu.DB, tenantid)

	return GetData, nil
}

func (menu *Menu) GetmenusByTenantId(websiteid int, tenantid string) ([]TblMenus, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return []TblMenus{}, AuthError
	}

	GetData, _ := menumodel.GetmenusByTenantId(websiteid,menu.DB, tenantid)

	return GetData, nil
}

func (menu *Menu) UpdateMenuItemOrder(menuitems []OrderItem, userid int, tenantid string) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}

	err := menumodel.UpdateMenuItemOrder(menu.DB, menuitems, userid, tenantid)

	if err != nil {

		return err
	}

	return nil
}
