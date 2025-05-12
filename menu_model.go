package menu

import (
	"time"

	"gorm.io/gorm"
)

type Filter struct {
	Keyword string
	Status  string
}

type TblMenus struct {
	Id          int
	Name        string
	Description string
	TenantId    string
	CreatedOn   time.Time
	CreatedBy   int
	IsDeleted   int
	DeletedOn   time.Time
	DeletedBy   int
	ModifiedOn  time.Time
	ModifiedBy  int
	DateString  string `gorm:"-"`
	ParentId    int
	UrlPath     string
	Slug        string
}

type MenuModel struct {
	Userid     int
	DataAccess int
}

type MenuCreate struct {
	Id          int
	MenuName    string
	MenuSlug    string
	Description string
	ParentId    int
	TenantId    string
	CreatedBy   int
	ModifiedBy  int
}
//Menu Listing
func (menu *MenuModel) MenuList(limit int, offset int, filter Filter, DB *gorm.DB, tenantid string) (menus []TblMenus, count int64, err error) {

	var menucount int64

	query := DB.Table("tbl_menus").Where("is_deleted = 0 and parent_id=0 and  tenant_id = ?", tenantid).Order("tbl_menus.created_on desc")

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(name)) like LOWER(TRIM(?))", "%"+filter.Keyword+"%")
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

//Create Menu
func (menu *MenuModel) CreateMenus(req *TblMenus, DB *gorm.DB) error {

	if err := DB.Table("tbl_menus").Create(&req).Error; err != nil {

		return err
	}
	return nil
}
