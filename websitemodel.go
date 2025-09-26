package menu

import (
	"time"

	"gorm.io/gorm"
)

type TblWebsite struct {
	Id            int       `gorm:"primaryKey;auto_increment;type:serial"`
	Name          string    `gorm:"type:character varying"`
	ChannelNames  string    `gorm:"type:character varying"`
	TemplateId    int       `gorm:"type:integer"`
	TenantId      string    `gorm:"type:character varying"`
	IsDeleted     int       `gorm:"type:integer"`
	DeletedOn     time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy     int       `gorm:"DEFAULT:NULL"`
	CreatedOn     time.Time `gorm:"type:timestamp without time zone"`
	CreatedBy     int       `gorm:"type:integer"`
	ModifiedOn    time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy    int       `gorm:"DEFAULT:NULL;type:integer"`
	Status        int       `gorm:"type:integer"`
	DateString    string    `gorm:"-:migration;<-:false"`
	TemplateName  string    `gorm:"column:template_name;->;<-:false"` // read-only from join
	TemplateImage string    `gorm:"column:template_image;->;<-:false"`
	CreatedDate   string    `gorm:"-:migration;<-:false"`
	Subdomain     string    `gorm:"-:migration;<-:false"`
}

// createwebsite
func (menu *MenuModel) CreateWebsite(website *TblWebsite, DB *gorm.DB) (TblWebsite, error) {

	if err := DB.Table("tbl_websites").Create(&website).Error; err != nil {

		return TblWebsite{}, err
	}
	return *website, nil
}

// websitelist
func (menu *MenuModel) WebsiteList(limit int, offset int, filter Filter, DB *gorm.DB, tenantid string) (website []TblWebsite, count int64, err error) {

	query := DB.Table("tbl_websites").
		Select("tbl_websites.*, tbl_go_templates.template_name, tbl_go_templates.template_image").
		Joins("left join tbl_go_templates on tbl_go_templates.id = tbl_websites.template_id").
		Where("tbl_websites.is_deleted = 0 AND tbl_websites.tenant_id = ?", tenantid).
		Order("tbl_websites.created_on desc")

	if filter.Keyword != "" {
		query = query.Where("LOWER(TRIM(tbl_websites.name)) like LOWER(TRIM(?))", "%"+filter.Keyword+"%")
	}
	if filter.ToDate != "" {
		query = query.Where("tbl_websites.modified_on >= ? AND tbl_websites.modified_on < ?", filter.ToDate+" 00:00:00", filter.ToDate+" 23:59:59")
	}
	if filter.Status != "" {
		switch filter.Status {
		case "publish":
			query = query.Where("tbl_websites.status = ?", 1)
		case "unpublish":
			query = query.Where("tbl_websites.status = ?", 0)
		}
	}

	err = query.Limit(limit).Offset(offset).Find(&website).Error
	if err != nil {
		return nil, 0, err
	}
	// Count total matching records without limit/offset
	err = DB.Table("tbl_websites").
		Where("tbl_websites.is_deleted = 0 AND tbl_websites.tenant_id = ?", tenantid).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return website, count, nil
}

//update website func//

func (menu *MenuModel) UpdateWebsite(website *TblWebsite, DB *gorm.DB) (TblWebsite, error) {

	if err := DB.Table("tbl_websites").Where("id = ? and  tenant_id = ?", website.Id, website.TenantId).UpdateColumns(map[string]interface{}{"name": website.Name, "channel_names": website.ChannelNames, "template_id": website.TemplateId, "status": website.Status, "modified_by": website.ModifiedBy, "modified_on": website.ModifiedOn}).Error; err != nil {

		return TblWebsite{}, err
	}

	return *website, nil
}

func (menu *MenuModel) GetWebsiteById(webid int, tenantid string, DB *gorm.DB) (website TblWebsite, err error) {

	if err := DB.Table("tbl_websites").Where("is_deleted = 0 and id=? and tenant_id=? ", webid, tenantid).Order("id asc").Find(&website).Error; err != nil {

		return TblWebsite{}, err
	}

	return website, nil
}

func (menu *MenuModel) GetWebsiteByName(name string, DB *gorm.DB) (website TblWebsite, err error) {

	if err := DB.Table("tbl_websites").Where("is_deleted = 0 and  name=? ", name).Order("name asc").Find(&website).Error; err != nil {

		return TblWebsite{}, err
	}

	return website, nil
}
func (menu *MenuModel) CheckSiteName(name string, webid int, DB *gorm.DB) error {

	var website TblWebsite

	if webid == 0 {
		if err := DB.Table("tbl_websites").Where("name = ? AND  is_deleted = 0 ", name).First(&website).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Table("tbl_websites").Where("name = ? AND id not in(?) and  is_deleted = 0 ", name, webid).First(&website).Error; err != nil {

			return err
		}
	}

	return nil

}

//Delete Website//

func (menu *MenuModel) DeleteWebsiteById(website *TblWebsite, DB *gorm.DB) error {

	if err := DB.Debug().Table("tbl_websites").Where("id=? and  tenant_id = ?", website.Id, website.TenantId).Updates(TblWebsite{IsDeleted: website.IsDeleted, DeletedOn: website.DeletedOn, DeletedBy: website.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}
