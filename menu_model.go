package menu

import (
	"time"

	"gorm.io/gorm"
)

type Filter struct {
	Keyword string
	Status  string
	ToDate  string
}

type TblMenus struct {
	Id              int
	Name            string
	Description     string
	TenantId        string
	CreatedOn       time.Time
	CreatedBy       int
	IsDeleted       int
	DeletedOn       time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy       int       `gorm:"DEFAULT:NULL"`
	ModifiedOn      time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy      int       `gorm:"DEFAULT:NULL"`
	DateString      string    `gorm:"-"`
	ParentId        int
	UrlPath         string
	SlugName        string
	Status          int
	Type            string
	TypeId          int
	MenuitemCount   int `gorm:"-"`
	Count           int `gorm:"-"`
	WebsiteId       int
	ListingsIds     string
	CategoryIds     string
	ImageName       string
	ImagePath       string
	OrderIndex      int
	MetaTitle       string
	MetaDescription string
	MetaKeywords    string
	SeparateWindow  int
}

type MenuModel struct {
	Userid     int
	DataAccess int
}

type MenuCreate struct {
	Id              int
	MenuName        string
	MenuSlug        string
	Description     string
	ParentId        int
	TenantId        string
	CreatedBy       int
	ModifiedBy      int
	Status          int
	UrlPath         string
	Type            string
	TypeId          int
	WebsiteId       int
	ListingsIds     string
	CategoryIds     string
	ImagePath       string
	ImageName       string
	MetaTitle       string
	MetaDescription string
	MetaKeywords    string
	SeperateWindow  int
	OrderIndex      int
}
type OrderItem struct {
	MenuItemID   int  `json:"menuitem_id"`
	OrderIndex   int  `json:"orderindex"`
	ParentMenuID int  `json:"parentmenu_id"`
	IsChild      bool `json:"is_child"`
}

// Menu Listing
func (menu *MenuModel) MenuList(limit int, offset int, filter Filter, DB *gorm.DB, tenantid string, websiteid int) (menus []TblMenus, count int64, err error) {

	var menucount int64

	query := DB.Table("tbl_menus").Where("is_deleted = 0 and parent_id=0 and website_id=? and tenant_id = ?", websiteid, tenantid).Order("tbl_menus.created_on desc")

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(name)) like LOWER(TRIM(?))", "%"+filter.Keyword+"%")
	}

	if filter.ToDate != "" {
		query = query.Where("tbl_menus.modified_on >= ? AND tbl_menus.modified_on < ?",
			filter.ToDate+" 00:00:00",
			filter.ToDate+" 23:59:59")
	}
	if filter.Status != "" {

		if filter.Status == "Active" {

			query = query.Where("tbl_menus.status=?", 1)
		}
		if filter.Status == "Inactive" {

			query = query.Where("tbl_menus.status=?", 0)
		}
	}
	if limit != 0 {

		query.Limit(limit).Offset(offset).Find(&menus)

		return menus, menucount, nil

	}

	query.Find(&menus).Count(&menucount)

	if query.Error != nil {

		return []TblMenus{}, 0, query.Error
	}

	return menus, menucount, nil
}

// Create Menu
func (menu *MenuModel) CreateMenus(req *TblMenus, DB *gorm.DB) (TblMenus, error) {

	if err := DB.Table("tbl_menus").Create(&req).Error; err != nil {

		return TblMenus{}, err
	}
	return *req, nil
}

// UpdateMenu
func (menu *MenuModel) UpdateMenu(menureq *TblMenus, DB *gorm.DB) (TblMenus, error) {

	if menureq.ParentId == 0 {

		if err := DB.Table("tbl_menus").Where("id = ? and  tenant_id = ?", menureq.Id, menureq.TenantId).UpdateColumns(map[string]interface{}{"name": menureq.Name, "slug_name": menureq.SlugName, "status": menureq.Status, "description": menureq.Description, "modified_by": menureq.ModifiedBy, "modified_on": menureq.ModifiedOn, "website_id": menureq.WebsiteId, "separate_window": menureq.SeparateWindow}).Error; err != nil {

			return TblMenus{}, err
		}
	} else {
		if err := DB.Table("tbl_menus").Where("id = ? and  tenant_id = ?", menureq.Id, menureq.TenantId).UpdateColumns(map[string]interface{}{"name": menureq.Name, "url_path": menureq.UrlPath, "parent_id": menureq.ParentId, "status": menureq.Status, "slug_name": menureq.SlugName, "modified_by": menureq.ModifiedBy, "modified_on": menureq.ModifiedOn, "type": menureq.Type, "type_id": menureq.TypeId, "website_id": menureq.WebsiteId, "image_name": menureq.ImageName, "image_path": menureq.ImagePath, "listings_ids": menureq.ListingsIds, "category_ids": menureq.CategoryIds, "meta_title": menureq.MetaTitle, "meta_description": menureq.MetaDescription, "meta_keywords": menureq.MetaKeywords, "separate_window": menureq.SeparateWindow}).Error; err != nil {

			return TblMenus{}, err
		}
	}

	return *menureq, nil
}

func (menu *MenuModel) GetMenuTree(menuid int, DB *gorm.DB, tenantid string) ([]TblMenus, error) {
	var menus []TblMenus
	err := DB.Raw(`
        WITH RECURSIVE me_tree AS (
            SELECT id,  name,
            slug_name,
            parent_id,
            created_on,
            modified_on,
            url_path,
            type,
            type_id,
            is_deleted,
            listings_ids,
            image_path,
            image_name,
            order_index,
            separate_window
            FROM tbl_menus
            WHERE id = ? and  tenant_id =?
            UNION ALL
            SELECT me.id, me.name,
            me.slug_name,
            me.parent_id,
            me.created_on,
            me.modified_on,
            me.url_path,
            me.type,
            me.type_id,
            me.is_deleted,
            me.listings_ids,
            me.image_path,
            me.image_name,
            me.order_index,
            me.separate_window
            FROM tbl_menus AS me
            JOIN me_tree ON me.parent_id = me_tree.id and  me.tenant_id =?
        )
        SELECT *
        FROM me_tree WHERE IS_DELETED = 0 order by order_index asc
    `, menuid, tenantid, tenantid).Scan(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

// DeleteMenuById
func (menu *MenuModel) DeleteMenuById(menureq *TblMenus, menuid []int, tenantid string, DB *gorm.DB) error {

	if err := DB.Table("tbl_menus").Where("id in(?) and  tenant_id = ?", menuid, tenantid).Updates(TblMenus{IsDeleted: menureq.IsDeleted, DeletedOn: menureq.DeletedOn, DeletedBy: menureq.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

// Check Menuname is already exists
func (menu *MenuModel) CheckMenuName(menureq TblMenus, menuid int, name string, parentid int, websiteid int, DB *gorm.DB, tenantid string) error {

	if menuid == 0 {

		if err := DB.Debug().Table("tbl_menus").Where("LOWER(TRIM(name))=LOWER(TRIM(?)) and is_deleted=0 and website_id=? and tenant_id = ? and parent_id=?", name, websiteid, tenantid, parentid).First(&menureq).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Table("tbl_menus").Where("LOWER(TRIM(name))=LOWER(TRIM(?)) and id not in (?) and is_deleted=0 and website_id=? and  tenant_id = ? and parent_id=?", name, menuid, websiteid, tenantid, parentid).First(&menureq).Error; err != nil {

			return err
		}
	}

	return nil
}

func (menu *MenuModel) MenuStatusChange(menureq TblMenus, DB *gorm.DB) error {

	if err := DB.Table("tbl_menus").Where("id=? and tenant_id=?", menureq.Id, menureq.TenantId).UpdateColumns(map[string]interface{}{"status": menureq.Status, "modified_by": menureq.ModifiedBy, "modified_on": menureq.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

func (menu *MenuModel) DeleteMenuItemById(menureq TblMenus, DB *gorm.DB) error {

	if err := DB.Table("tbl_menus").Where("id=?  and  tenant_id = ?", menureq.Id, menureq.TenantId).Updates(TblMenus{IsDeleted: menureq.IsDeleted, DeletedOn: menureq.DeletedOn, DeletedBy: menureq.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

func (menu *MenuModel) GetMenuById(menuid int, DB *gorm.DB, tenantid string) (TblMenus, error) {

	var menudet TblMenus

	if err := DB.Table("tbl_menus").Where("id=? and tenant_id=? and is_deleted=0", menuid, tenantid).First(&menudet).Error; err != nil {

		return TblMenus{}, err
	}
	return menudet, nil
}

func (menu *MenuModel) GetMenuBySlug(menuslug string, DB *gorm.DB, tenantid string) (TblMenus, error) {

	var menudet TblMenus

	if err := DB.Table("tbl_menus").Where("slug_name=? and tenant_id=? and is_deleted=0", menuslug, tenantid).First(&menudet).Error; err != nil {

		return TblMenus{}, err
	}

	return menudet, nil
}

func (menu *MenuModel) GetMenuBySlugName(menuslug string, websiteid int, DB *gorm.DB, tenantid string) (TblMenus, error) {

	var menudet TblMenus

	if err := DB.Table("tbl_menus").Where("slug_name=? and tenant_id=? and is_deleted=0 and website_id=?", menuslug, tenantid, websiteid).First(&menudet).Error; err != nil {

		return TblMenus{}, err
	}

	return menudet, nil
}

func (menu *MenuModel) GetMenuByUrlPath(urlpath string, websiteid int, DB *gorm.DB, tenantid string) (TblMenus, error) {

	var menudet TblMenus

	if err := DB.Debug().Table("tbl_menus").Where("url_path=? and tenant_id=? and is_deleted=0 and website_id=?", urlpath, tenantid, websiteid).First(&menudet).Error; err != nil {

		return TblMenus{}, err
	}

	return menudet, nil
}

func (menu *MenuModel) GetmenusByTenantId(websiteid int, DB *gorm.DB, tenantid string) ([]TblMenus, error) {

	var menudet []TblMenus

	if err := DB.Table("tbl_menus").Where("tenant_id=? and website_id=? and is_deleted=0 and status=1", tenantid, websiteid).Order("order_index asc").Find(&menudet).Error; err != nil {

		return []TblMenus{}, err
	}

	return menudet, nil
}
func (menu *MenuModel) UpdateMenuItemOrder(DB *gorm.DB, menuitems []OrderItem, userid int, tenantid string) error {
	ModifiedOn, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	for _, item := range menuitems {

		if err := DB.Table("tbl_menus").Where("id=? and tenant_id=? and is_deleted=0", item.MenuItemID, tenantid).UpdateColumns(map[string]interface{}{"order_index": item.OrderIndex, "parent_id": item.ParentMenuID, "modified_by": userid, "modified_on": ModifiedOn}).Error; err != nil {

			return err
		}
	}
	return nil
}

func (menu *MenuModel) GetDirectSubMenusByParentID(parentid int, tenantid string, DB *gorm.DB) (menus []TblMenus, err error) {

	if err := DB.Table("tbl_menus").Where("is_deleted=0 and parent_id=? and tenant_id=?", parentid, tenantid).Find(&menus).Error; err != nil {

		return []TblMenus{}, err

	}

	return menus, nil
}

func (menu *MenuModel) UpdateMenuOrderIndex(menus *TblMenus, menuid int, parent_id int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_menus").Where("id=? and parent_id=? and tenant_id=?", menuid, parent_id, menus.TenantId).UpdateColumns(map[string]interface{}{"order_index": menus.OrderIndex}).Error; err != nil {

		return err
	}

	return nil
}
