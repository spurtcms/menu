package menu

import (
	"time"

	"gorm.io/gorm"
)

type TblTemplatePages struct {
	Id              int       `gorm:"primaryKey;auto_increment;type:serial"`
	Name            string    `gorm:"type:character varying"`
	Slug            string    `gorm:"type:character varying"`
	PageDescription string    `gorm:"type:character varying"`
	TenantId        string    `gorm:"type:character varying"`
	IsDeleted       int       `gorm:"type:integer"`
	DeletedOn       time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy       int       `gorm:"DEFAULT:NULL"`
	CreatedOn       time.Time `gorm:"type:timestamp without time zone"`
	CreatedBy       int       `gorm:"DEFAULT:NULL"`
	ModifiedOn      time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy      int       `gorm:"DEFAULT:NULL"`
	CreatedDate     string    `gorm:"-:migration;<-:false"`
	ModifiedDate    string    `gorm:"-:migration;<-:false"`
	Status          int       `gorm:"type:integer"`
	MetaTitle       string    `gorm:"type:character varying"`
	MetaDescription string    `gorm:"type:character varying"`
	MetaKeywords    string    `gorm:"type:character varying"`
	MetaSlug        string    `gorm:"type:character varying"`
	WebsiteId       int       `gorm:"type:integer"`
	MenuNames       string    `gorm:"-"`
	PageType        string    `gorm:"type:character varying"`
	CustomPagePath  string    `gorm:"type:character varying"`
	ParentId        int       `gorm:"type:integer"`
	OrderIndex      int       `gorm:"type:integer"`
}

// Create Page
func (menu *MenuModel) CreateTemplatePage(db *gorm.DB, page *TblTemplatePages) (TblTemplatePages, error) {

	if err := db.Table("tbl_template_pages").Create(&page).Error; err != nil {

		return TblTemplatePages{}, err
	}
	return *page, nil

}

// PageList
func (menu *MenuModel) TemplatePageList(limit int, offset int, filter Filter, DB *gorm.DB, Tenantid string, websiteid int) (pages []TblTemplatePages, count int64, err error) {

	var pagecount int64

	query := DB.Table("tbl_template_pages").Where("is_deleted = 0 and website_id=? and tenant_id = ?", websiteid, Tenantid).Order("tbl_template_pages.order_index asc")

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(name)) like LOWER(TRIM(?))", "%"+filter.Keyword+"%")
	}

	if filter.ToDate != "" {
		query = query.Where("tbl_template_pages.modified_on >= ? AND tbl_template_pages.modified_on < ?",
			filter.ToDate+" 00:00:00",
			filter.ToDate+" 23:59:59")
	}
	if filter.Status != "" {

		if filter.Status == "Active" {

			query = query.Where("tbl_template_pages.status=?", 1)
		}
		if filter.Status == "Inactive" {

			query = query.Where("tbl_template_pages.status=?", 0)
		}
	}
	if limit != 0 {

		query.Limit(limit).Offset(offset).Find(&pages)

		return pages, pagecount, nil

	}

	query.Find(&pages).Count(&pagecount)

	if query.Error != nil {

		return []TblTemplatePages{}, 0, query.Error
	}

	return pages, pagecount, nil

}

// Get PageById
func (menu *MenuModel) GetPageById(DB *gorm.DB, pageid int, tenantid string) (page TblTemplatePages, err error) {

	if err := DB.Table("tbl_template_pages").Where("is_deleted = 0 and id=? and tenant_id=?", pageid, tenantid).Order("id asc").Find(&page).Error; err != nil {

		return TblTemplatePages{}, err
	}

	return page, nil

}

// Update Page
func (menu *MenuModel) UpdateTemplatePage(db *gorm.DB, page *TblTemplatePages) (TblTemplatePages, error) {

	if err := db.Table("tbl_template_pages").Where("id = ? and tenant_id=?", page.Id, page.TenantId).Updates(page).Error; err != nil {
		return TblTemplatePages{}, err
	}
	return *page, nil

}

// Status change func
func (menu *MenuModel) PageStatusChange(page TblTemplatePages, DB *gorm.DB) error {

	if err := DB.Table("tbl_template_pages").Where("id=? and tenant_id=?", page.Id, page.TenantId).UpdateColumns(map[string]interface{}{"status": page.Status, "modified_by": page.ModifiedBy, "modified_on": page.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

// DeletePageById
func (menu *MenuModel) DeletePageById(page *TblTemplatePages, tenantid string, DB *gorm.DB) error {

	if err := DB.Debug().Table("tbl_template_pages").Where("id=? and  tenant_id = ?", page.Id, tenantid).Updates(TblTemplatePages{IsDeleted: page.IsDeleted, DeletedOn: page.DeletedOn, DeletedBy: page.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

// Get PageBySlug
func (menu *MenuModel) GetPageBySlug(DB *gorm.DB, pageslug string, tenantid string) (page TblTemplatePages, err error) {

	if err := DB.Table("tbl_template_pages").Where("is_deleted = 0 and slug=? and tenant_id=?", pageslug, tenantid).Order("id asc").Find(&page).Error; err != nil {

		return TblTemplatePages{}, err
	}

	return page, nil

}

func (menu *MenuModel) GetMenusByPageId(DB *gorm.DB, pageid int, tenantid string) (menus TblMenus, err error) {

	if err := DB.Table("tbl_menus").Where("is_deleted = 0 and type=? and tenant_id=? and type_id=?", "pages", tenantid, pageid).Order("id asc").Find(&menus).Error; err != nil {

		return TblMenus{}, err
	}

	return menus, nil
}

// Check Menuname is already exists
func (menu *MenuModel) CheckPageNameIsExits(pagererq TblTemplatePages, menuid int, name string, websiteid int, DB *gorm.DB, tenantid string) error {

	if menuid == 0 {

		if err := DB.Debug().Table("tbl_template_pages").Where("LOWER(TRIM(name))=LOWER(TRIM(?)) and is_deleted=0 and website_id=? and tenant_id = ? ", name, websiteid, tenantid).First(&pagererq).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Table("tbl_template_pages").Where("LOWER(TRIM(name))=LOWER(TRIM(?)) and id not in (?) and is_deleted=0 and website_id=? and  tenant_id = ? ", name, menuid, websiteid, tenantid).First(&pagererq).Error; err != nil {

			return err
		}
	}

	return nil
}

func (menu *MenuModel) UpdatePagesOrder(DB *gorm.DB, pages []OrderItem, userid int, tenantid string) error {
	ModifiedOn, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	for _, item := range pages {

		if err := DB.Debug().Table("tbl_template_pages").Where("id=? and tenant_id=? and is_deleted=0", item.MenuItemID, tenantid).UpdateColumns(map[string]interface{}{"order_index": item.OrderIndex, "parent_id": item.ParentMenuID, "modified_by": userid, "modified_on": ModifiedOn}).Error; err != nil {

			return err
		}
	}
	return nil
}
