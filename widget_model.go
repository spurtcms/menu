package menu

import (
	"time"

	"github.com/spurtcms/channels"
	"github.com/spurtcms/listing"
	"gorm.io/gorm"
)

type TblWidgets struct {
	Id              int       `gorm:"primaryKey;auto_increment;type:serial"`
	Title           string    `gorm:"type:character varying"`
	LongTitle       string    `gorm:"type:character varying"`
	Slug            string    `gorm:"type:character varying"`
	Position        string    `gorm:"type:character varying"`
	SortOrder       int       `gorm:"type:integer"`
	WidgetType      string    `gorm:"type:character varying"`
	TenantId        string    `gorm:"type:character varying"`
	WebsiteId       int       `gorm:"type:integer"`
	Status          int       `gorm:"type:integer;DEFAULT:1"`
	MetaTitle       string    `gorm:"type:character varying"`
	MetaDescription string    `gorm:"type:character varying"`
	MetaKeywords    string    `gorm:"type:character varying"`
	IsDeleted       int       `gorm:"type:integer;DEFAULT:0"`
	DeletedOn       time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy       int       `gorm:"DEFAULT:NULL"`
	CreatedOn       time.Time `gorm:"type:timestamp without time zone"`
	CreatedBy       int       `gorm:"type:integer"`
	ModifiedOn      time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy      int       `gorm:"DEFAULT:NULL;type:integer"`
	CreatedDate     string    `gorm:"-:migration;<-:false"`
	ModifiedDate    string    `gorm:"-:migration;<-:false"`
	ProductIds      string    `gorm:"-:migration;<-:false"`

	EntriesData           []channels.Tblchannelentries `gorm:"-"`
	ListingData           []listing.TblListing         `gorm:"-"`
	CategoryBaseEntryData []channels.Tblchannelentries `gorm:"-"`

	WidgetTitle string `gorm:"-:migration;<-:false"`
	WidgetId    int    `gorm:"-:migration;<-:false"`
}

type TblWidgetProducts struct {
	Id        int       `gorm:"primaryKey;auto_increment;type:serial"`
	WidgetId  int       `gorm:"type:integer"`
	ProductId int       `gorm:"type:integer"`
	TenantId  string    `gorm:"type:character varying"`
	CreatedOn time.Time `gorm:"type:timestamp without time zone"`
	CreatedBy int       `gorm:"type:integer"`
}

// WidgetList
func (menu *MenuModel) WidgetList(limit int, offset int, filter Filter, DB *gorm.DB, Tenantid string, websiteid int) (widgets []TblWidgets, count int64, err error) {

	var widgetcount int64

	query := DB.Table("tbl_widgets").Where("is_deleted = 0 and website_id=? and tenant_id = ?", websiteid, Tenantid).Order("tbl_widgets.created_on desc")

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(title)) like LOWER(TRIM(?))", "%"+filter.Keyword+"%")
	}

	if filter.ToDate != "" {
		query = query.Where("tbl_widgets.modified_on >= ? AND tbl_widgets.modified_on < ?",
			filter.ToDate+" 00:00:00",
			filter.ToDate+" 23:59:59")
	}
	if filter.Status != "" {

		if filter.Status == "Active" {

			query = query.Where("tbl_widgets.status=?", 1)
		}
		if filter.Status == "Inactive" {

			query = query.Where("tbl_widgets.status=?", 0)
		}
	}
	if limit != 0 {

		query.Limit(limit).Offset(offset).Find(&widgets)

		return widgets, widgetcount, nil

	}

	query.Find(&widgets).Count(&widgetcount)

	if query.Error != nil {

		return []TblWidgets{}, 0, query.Error
	}

	return widgets, widgetcount, nil

}

//create Widget

func (menu *MenuModel) CreateWidget(db *gorm.DB, widget *TblWidgets) (TblWidgets, error) {

	if err := db.Table("tbl_widgets").Create(&widget).Error; err != nil {

		return TblWidgets{}, err
	}
	return *widget, nil

}

func (menu *MenuModel) InsertWidgetProductIds(db *gorm.DB, widget *TblWidgetProducts) error {

	if err := db.Table("tbl_widget_products").Create(&widget).Error; err != nil {

		return err
	}

	return nil
}

// Get WidgetById
func (menu *MenuModel) GetWidgetById(DB *gorm.DB, widgetid int, tenantid string) (widget TblWidgets, products []TblWidgetProducts, err error) {

	if err := DB.Debug().
		Table("tbl_widgets").
		Where("is_deleted = 0 AND id = ? AND tenant_id = ?", widgetid, tenantid).
		Take(&widget).Error; err != nil {
		return TblWidgets{}, nil, err
	}

	if err := DB.Debug().
		Table("tbl_widget_products").
		Where("widget_id = ? and tenant_id=?", widgetid, tenantid).
		Find(&products).Error; err != nil {
		return widget, nil, err
	}

	return widget, products, nil
}

// Update Widget
func (menu *MenuModel) UpdateWidget(db *gorm.DB, widget *TblWidgets) (TblWidgets, error) {

	if err := db.Table("tbl_widgets").Where("id = ? and tenant_id=?", widget.Id, widget.TenantId).UpdateColumns(map[string]interface{}{"status": widget.Status, "modified_by": widget.ModifiedBy, "modified_on": widget.ModifiedOn, "title": widget.Title, "long_title": widget.LongTitle, "position": widget.Position, "sort_order": widget.SortOrder, "meta_title": widget.MetaTitle, "meta_description": widget.MetaDescription, "meta_keywords": widget.MetaKeywords, "slug": widget.Slug, "widget_type": widget.WidgetType, "website_id": widget.WebsiteId}).Error; err != nil {
		return TblWidgets{}, err
	}
	return *widget, nil

}

// Status change func
func (menu *MenuModel) WidgetStatusChange(page TblWidgets, DB *gorm.DB) error {

	if err := DB.Table("tbl_widgets").Where("id=? and tenant_id=?", page.Id, page.TenantId).UpdateColumns(map[string]interface{}{"status": page.Status, "modified_by": page.ModifiedBy, "modified_on": page.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

// DeleteWidgetById
func (menu *MenuModel) DeleteWidgetById(page *TblWidgets, tenantid string, DB *gorm.DB) error {

	if err := DB.Debug().Table("tbl_widgets").Where("id=? and  tenant_id = ?", page.Id, tenantid).Updates(TblWidgets{IsDeleted: page.IsDeleted, DeletedOn: page.DeletedOn, DeletedBy: page.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

// Get PageBySlug
func (menu *MenuModel) GetWidgetBySlug(DB *gorm.DB, pageslug string, tenantid string) (page TblWidgets, err error) {

	if err := DB.Table("tbl_widgets").Where("is_deleted = 0 and slug=? and tenant_id=?", pageslug, tenantid).Order("id asc").Find(&page).Error; err != nil {

		return TblWidgets{}, err
	}

	return page, nil

}

//Delete ProductIds

func (menu *MenuModel) DeleteProductIds(DB *gorm.DB, widgetid int, tenantid string) error {

	if err := DB.Table("tbl_widget_products").
		Where("widget_id = ? AND tenant_id = ?", widgetid, tenantid).
		Delete(nil).Error; err != nil {

		return err
	}
	return nil
}

func (menu *MenuModel) FetchBasicWidgetList(DB *gorm.DB, tenantID string, websiteID int) ([]TblWidgets, error) {
	var widgets []TblWidgets
	err := DB.Debug().Table("tbl_widgets").Where("is_deleted = 0 AND tenant_id = ? AND website_id = ? AND status = 1", tenantID, websiteID).
		Find(&widgets).Error
	return widgets, err
}

func (menu *MenuModel) FetchWidgetEntries(DB *gorm.DB, widgetID int) ([]channels.Tblchannelentries, error) {
	var entries []channels.Tblchannelentries
	err := DB.Table("tbl_widget_products AS wp").
		Select("ce.*,c.slug_name as channel_name").
		Joins("JOIN tbl_channel_entries AS ce ON wp.product_id = ce.id").Joins("left join tbl_channels as c on c.id =ce.channel_id").
		Where("wp.widget_id = ? and ce.is_deleted=0", widgetID).Limit(6).
		Find(&entries).Error
	return entries, err
}
func (menu *MenuModel) FetchWidgetByCategoriesEntries(DB *gorm.DB, widgetID int) ([]channels.Tblchannelentries, error) {
	var entries []channels.Tblchannelentries
	err := DB.
		Table("tbl_widget_products AS wp").
		Select("ce.*, c.slug_name as channel_name").
		Joins("JOIN tbl_channel_entries AS ce ON ?::text = ANY(string_to_array(ce.categories_id, ','))", gorm.Expr("CAST(wp.product_id AS text)")).
		Joins("LEFT JOIN tbl_channels AS c ON c.id = ce.channel_id").
		Where("wp.widget_id = ? and ce.is_deleted=0", widgetID).Limit(6).
		Find(&entries).Error

	return entries, err
}
func (menu *MenuModel) FetchWidgetListings(DB *gorm.DB, widgetID int) ([]listing.TblListing, error) {
	var listings []listing.TblListing
	err := DB.Debug().Table("tbl_widget_products AS wp").
		Select("l.*, ce.slug as entry_slug").
		Joins("JOIN tbl_listings AS l ON wp.product_id = l.id").Joins("left join tbl_channel_entries as ce on ce.id =l.entry_id").
		Where("wp.widget_id = ? and l.is_deleted=0", widgetID).Limit(6).
		Find(&listings).Error
	return listings, err
}
