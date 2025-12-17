package menu

import (
	"fmt"
	"strings"
	"time"
)

// PageList
func (menu *Menu) GetTemplatePageList(limit int, offset int, filter Filter, tenantid string, websiteid int) ([]TblTemplatePages, int, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return []TblTemplatePages{}, 0, AuthError
	}

	menumodel.DataAccess = menu.DataAccess
	menumodel.Userid = menu.UserId

	_, totalcount, _ := menumodel.TemplatePageList(0, 0, filter, menu.DB, tenantid, websiteid)

	pagelist, _, cerr := menumodel.TemplatePageList(limit, offset, filter, menu.DB, tenantid, websiteid)

	if cerr != nil {

		return []TblTemplatePages{}, 0, cerr
	}

	return pagelist, int(totalcount), nil
}

// CreatePage
func (menu *Menu) CreateTemplatePage(page *TblTemplatePages) (TblTemplatePages, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblTemplatePages{}, AuthError
	}
	page.Slug = strings.ToLower(strings.ReplaceAll(page.Name, " ", "-"))
	page.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	pagedetail, err := menumodel.CreateTemplatePage(menu.DB, page)

	if err != nil {

		return TblTemplatePages{}, err

	}
	return pagedetail, nil

}

//GetPagebyId

func (menu *Menu) GetPageById(pageid int, tenantid string) (TblTemplatePages, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblTemplatePages{}, AuthError
	}
	pagedetail, err := menumodel.GetPageById(menu.DB, pageid, tenantid)

	if err != nil {

		return TblTemplatePages{}, err

	}
	return pagedetail, nil
}

// UpdatePage
func (menu *Menu) EditTemplatePage(page *TblTemplatePages) (TblTemplatePages, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblTemplatePages{}, AuthError
	}

	page.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	pagedetail, err := menumodel.UpdateTemplatePage(menu.DB, page)

	if err != nil {

		return TblTemplatePages{}, err

	}
	return pagedetail, nil

}

//Delete page

func (menu *Menu) DeletePage(pageid int, modifiedby int, tenantid string) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}

	var pagedet TblTemplatePages

	pagedet.DeletedBy = modifiedby

	pagedet.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	pagedet.IsDeleted = 1

	pagedet.Id = pageid

	err := menumodel.DeletePageById(&pagedet, tenantid, menu.DB)

	if err != nil {

		return err
	}

	return nil

}

// Status chage
func (menu *Menu) PageStatusChange(pageid int, status int, userid int, tenantid string) (bool, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return false, AuthError
	}

	var pagedet TblTemplatePages

	pagedet.ModifiedBy = userid

	pagedet.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	pagedet.Status = status

	pagedet.TenantId = tenantid

	pagedet.Id = pageid

	err := menumodel.PageStatusChange(pagedet, menu.DB)

	if err != nil {

		return false, err

	}
	return true, nil
}

// GetpagebySlug
func (menu *Menu) GetPageBySlug(slug string, tenantid string) (TblTemplatePages, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblTemplatePages{}, AuthError
	}
	pagedetail, err := menumodel.GetPageBySlug(menu.DB, slug, tenantid)

	if err != nil {

		return TblTemplatePages{}, err

	}
	return pagedetail, nil

}

func (menu *Menu) GetMenusByPageId(pageid int, tenantid string) (TblMenus, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblMenus{}, AuthError
	}

	menus, err := menumodel.GetMenusByPageId(menu.DB, pageid, tenantid)

	if err != nil {

		fmt.Println(err)
	}

	return menus, nil
}

// Check Menuname is already exists
func (menu *Menu) CheckPageNameIsExits(id int, name string, websiteid int, tenantid string) (bool, error) {

	var pages TblTemplatePages

	err := menumodel.CheckPageNameIsExits(pages, id, name, websiteid, menu.DB, tenantid)

	if err != nil {

		return false, err

	}

	return true, nil
}
func (menu *Menu) UpdatePagesOrder(pages []OrderItem, userid int, tenantid string) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}

	err := menumodel.UpdatePagesOrder(menu.DB, pages, userid, tenantid)

	if err != nil {

		return err
	}

	return nil
}
