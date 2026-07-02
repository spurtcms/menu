package menu

import (
	"fmt"
	"html/template"
	"time"

	"gorm.io/gorm"
)

type TblTemplatePages struct {
	Id              int           `gorm:"primaryKey;auto_increment;type:serial"`
	Name            string        `gorm:"type:character varying"`
	Slug            string        `gorm:"type:character varying"`
	PageDescription string        `gorm:"type:character varying"`
	TenantId        string        `gorm:"type:character varying"`
	IsDeleted       int           `gorm:"type:integer"`
	DeletedOn       time.Time     `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy       int           `gorm:"DEFAULT:NULL"`
	CreatedOn       time.Time     `gorm:"type:timestamp without time zone"`
	CreatedBy       int           `gorm:"DEFAULT:NULL"`
	ModifiedOn      time.Time     `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy      int           `gorm:"DEFAULT:NULL"`
	CreatedDate     string        `gorm:"-:migration;<-:false"`
	ModifiedDate    string        `gorm:"-:migration;<-:false"`
	Status          int           `gorm:"type:integer"`
	MetaTitle       string        `gorm:"type:character varying"`
	MetaDescription string        `gorm:"type:character varying"`
	MetaKeywords    string        `gorm:"type:character varying"`
	MetaSlug        string        `gorm:"type:character varying"`
	OgImage         string        `gorm:"type:character varying"`
	WebsiteId       int           `gorm:"type:integer"`
	MenuNames       string        `gorm:"-"`
	PageType        string        `gorm:"type:character varying"`
	CustomPagePath  string        `gorm:"type:character varying"`
	ParentId        int           `gorm:"type:integer"`
	OrderIndex      int           `gorm:"type:integer"`
	HtmlDescription template.HTML `gorm:"-"`
	CloneCount      int           `gorm:"type:integer"`
	StructureId     int           `gorm:"type:integer"`
	GroupId         int           `gorm:"type:integer"`
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

	query := DB.Table("tbl_template_pages").Where("is_deleted = 0  and tenant_id = ?", Tenantid).Order("tbl_template_pages.order_index asc")

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

	if filter.PageId != 0 {
		query = query.Where("parent_id=?", filter.PageId)
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
func (menu *MenuModel) DeletePageById(page *TblTemplatePages, individualid []int, tenantid string, DB *gorm.DB) error {

	// if err := DB.Debug().Table("tbl_template_pages").Where("id=? and  tenant_id = ?", page.Id, tenantid).Updates(TblTemplatePages{IsDeleted: page.IsDeleted, DeletedOn: page.DeletedOn, DeletedBy: page.DeletedBy}).Error; err != nil {

	// 	return err

	// }

	if err := DB.Table("tbl_template_pages").Where("id in(?) and  tenant_id = ?", individualid, tenantid).Updates(TblTemplatePages{IsDeleted: page.IsDeleted, DeletedOn: page.DeletedOn, DeletedBy: page.DeletedBy}).Error; err != nil {

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

func (menu *MenuModel) GetPageTree(pageid int, DB *gorm.DB, tenantid string) ([]TblTemplatePages, error) {
	var pages []TblTemplatePages
	err := DB.Debug().Raw(`
		WITH RECURSIVE cat_tree AS (
			SELECT id, 	name,
			slug,
			parent_id,
			created_on,
			modified_on,
			is_deleted
			FROM tbl_template_pages
			WHERE id = ? and  tenant_id =?
			UNION ALL
			SELECT cat.id, cat.name,
			cat.slug,
			cat.parent_id,
			cat.created_on,
			cat.modified_on,
			cat.is_deleted
			FROM tbl_template_pages AS cat
			JOIN cat_tree ON cat.parent_id = cat_tree.id and  cat.tenant_id =?
		)
		SELECT *
		FROM cat_tree WHERE IS_DELETED = 0 order by id desc
	`, pageid, tenantid, tenantid).Scan(&pages).Error
	if err != nil {
		return nil, err
	}

	return pages, nil
}

func (menu *MenuModel) CloneCountUpdate(pageinfo TblTemplatePages, DB *gorm.DB) (Error error) {

	if err := DB.Table("tbl_template_pages").Where("id=?", pageinfo.Id).Updates(map[string]interface{}{"clone_count": pageinfo.CloneCount}).Error; err != nil {

		return err

	}

	return nil
}

//update Page OrderIndex query//

func (menu *MenuModel) UpdatePageOrderIndex(pageinfo *TblTemplatePages, DB *gorm.DB) error {

	if err := DB.Table("tbl_template_pages").Where("id=? and tenant_id=?", pageinfo.Id, pageinfo.TenantId).UpdateColumns(map[string]interface{}{"order_index": pageinfo.OrderIndex}).Error; err != nil {

		return err
	}

	return nil

}

// new page group and structure functions

type TblStructures struct {
	Id                   int       `gorm:"primaryKey;auto_increment;type:serial"`
	StructureName        string    `gorm:"type:character varying"`
	StructureSlug        string    `gorm:"type:character varying"`
	StructureDescription string    `gorm:"type:character varying"`
	TenantId             string    `gorm:"type:character varying"`
	CreatedOn            time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	CreatedBy            string    `gorm:"type:integer"`
	ModifiedOn           time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	IsDeleted            int       `gorm:"type:integer;DEFAULT:0"`
	ModifiedBy           int       `gorm:"type:integer"`
}

type StructureListResponse struct {
	ID int `json:"id"`

	StructureName string `json:"structure_name"`

	StructureSlug string `json:"structure_slug"`

	StructureDescription string `json:"structure_description"`

	TenantId string `json:"tenant_id"`

	PageCount int `json:"page_count"`

	PageGroupCount int `json:"page_group_count"`
}

type TblPageGroup struct {
	Id          int       `gorm:"primaryKey;auto_increment;type:serial"`
	GroupName   string    `gorm:"type:character varying"`
	GroupSlug   string    `gorm:"type:character varying"`
	TenantId    string    `gorm:"type:character varying"`
	StructureId int       `gorm:"type:integer"`
	CreatedOn   time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	CreatedBy   string    `gorm:"type:integer"`
	ModifiedOn  time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	IsDeleted   int       `gorm:"type:integer;DEFAULT:0"`
	ModifiedBy  int       `gorm:"type:integer"`
}

type PageGroupResponse struct {
	Id int `json:"id"`

	GroupName string `json:"group_name"`

	GroupSlug string `json:"group_slug"`

	Pages []PageTreeNode `json:"pages"`
}
type PageTreeNode struct {
	TblTemplatePages
	Children []TblTemplatePages
}

type StructureDetailsResponse struct {

	// structure details
	TblStructures TblStructures `json:"structure"`

	// top-level pages (parent_id=0) with children nested
	Pages []PageTreeNode `json:"pages"`

	// page groups with pages
	PageGroups []PageGroupResponse `json:"page_groups"`
}

type TblTemplatePagesResponce struct {
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
	StructureId     int       `gorm:"type:integer"`
	PagegroupId     int       `gorm:"type:integer"`

	CloneCount int `gorm:"type:integer"`
}

// page models

func (menu *MenuModel) Addpagegroupdata(group *TblPageGroup, DB *gorm.DB) (err error) {

	err1 := DB.Table("tbl_page_groups").Create(group).Error

	if err1 != nil {
		return err1
	}

	return nil

}

func (menu *MenuModel) GetStructureDetailsBasedonId(structureid int, DB *gorm.DB) (StructureDetails TblStructures, err error) {

	var structure TblStructures

	err1 := DB.Table("tbl_structures").Where("id = ?", structureid).Find(&structure).Error

	if err1 != nil {
		return structure, err1
	}

	return structure, nil

}

func (menu *MenuModel) Addstructuredata(structure TblStructures, DB *gorm.DB) (err error) {

	err1 := DB.Table("tbl_structures").Create(&structure).Error

	if err != nil {
		return err1
	}

	return nil

}

func (menu *MenuModel) GetStructureDataBasedOnTenant(Tenantid string, DB *gorm.DB) ([]StructureListResponse, error) {

	fmt.Println("GetStructureDataBasedOnTenantGetStructureDataBasedOnTenant tenant id", Tenantid)

	var structures []StructureListResponse

	err := DB.Table("tbl_structures s").
		Select(`
            s.id,
            s.structure_name,
            s.structure_slug,
            s.structure_description,
            s.tenant_id,
 
            (
                SELECT COUNT(*)
                FROM tbl_template_pages p
                WHERE p.structure_id = s.id
            ) as page_count,
 
            (
                SELECT COUNT(*)
                FROM tbl_page_groups g
                WHERE g.structure_id = s.id
            ) as page_group_count
        `).
		Where("s.tenant_id = ?", Tenantid).
		Scan(&structures).Error

	if err != nil {
		return nil, err
	}

	return structures, nil
}

func (menu *MenuModel) GetStructureDetails(structure_slug string, DB *gorm.DB) (StructureDetailsResponse, error) {

	var response StructureDetailsResponse

	// get structure details

	var structure TblStructures

	err := DB.Debug().
		Table("tbl_structures").
		Where(
			"structure_slug = ?",
			structure_slug,
		).
		First(&structure).Error

	if err != nil {
		return response, err
	}

	// assign structure data

	response.TblStructures = structure

	// get top-level direct pages (no group, no parent)

	var topPages []TblTemplatePages

	err = DB.Debug().
		Table("tbl_template_pages").
		Where(
			"structure_id = ? AND (group_id = 0 OR group_id IS NULL) AND parent_id = 0",
			structure.Id,
		).
		Find(&topPages).Error

	if err != nil {
		return response, err
	}

	for _, p := range topPages {
		var children []TblTemplatePages
		DB.Debug().
			Table("tbl_template_pages").
			Where("parent_id = ?", p.Id).
			Find(&children)
		response.Pages = append(response.Pages, PageTreeNode{
			TblTemplatePages: p,
			Children:         children,
		})
	}

	// get page groups

	var groups []TblPageGroup

	err = DB.Debug().
		Table("tbl_page_groups").
		Where(
			"structure_id = ?",
			structure.Id,
		).
		Find(&groups).Error

	if err != nil {
		return response, err
	}

	for _, group := range groups {

		var topGroupPages []TblTemplatePages

		DB.Debug().
			Table("tbl_template_pages").
			Where(
				"group_id = ? AND parent_id = 0",
				group.Id,
			).
			Find(&topGroupPages)

		var groupPageNodes []PageTreeNode
		for _, p := range topGroupPages {
			var children []TblTemplatePages
			DB.Debug().
				Table("tbl_template_pages").
				Where("parent_id = ?", p.Id).
				Find(&children)
			groupPageNodes = append(groupPageNodes, PageTreeNode{
				TblTemplatePages: p,
				Children:         children,
			})
		}

		response.PageGroups = append(
			response.PageGroups,

			PageGroupResponse{

				Id: group.Id,

				GroupName: group.GroupName,

				GroupSlug: group.GroupSlug,

				Pages: groupPageNodes,
			},
		)

	}

	return response, nil
}
