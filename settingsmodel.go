package menu

import (
	"errors"
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TblGoTemplateSettings struct {
	Id              int
	SiteName        string
	SiteLogo        string
	SiteLogoPath    string
	SiteFavIcon     string
	SiteFavIconPath string
	WebsiteUrl      string
	TenantId        string
	WebsiteId       int
	TemplateType    datatypes.JSON
	SocialMediaLink datatypes.JSON
	SocialLinks     []SocialLinks `gorm:"-"`
	HeaderThame     string
	TemplateID      string
}

type SocialLinks struct {
	Type      string
	IsActive  int
	SocialUrl string
}

func (menu *MenuModel) SettingDetail(tenantid string, websiteid int, DB *gorm.DB) (setting TblGoTemplateSettings, err error) {

	var SettingsDetail TblGoTemplateSettings

	if err := DB.Table("tbl_go_template_settings").Where("tenant_id = ? and website_id=?", tenantid, websiteid).First(&SettingsDetail).Error; err != nil {

		return TblGoTemplateSettings{}, err
	}

	return SettingsDetail, nil
}

func (menu *MenuModel) SettingDetailBasedONTemp(TemplateID string, tenantid string, websiteid int, DB *gorm.DB) (setting TblGoTemplateSettings, err error) {

	var SettingsDetail TblGoTemplateSettings

	if err := DB.Table("tbl_go_template_settings").Where("tenant_id = ? and website_id=? and template_id =?", tenantid, websiteid, TemplateID).First(&SettingsDetail).Error; err != nil {

		return TblGoTemplateSettings{}, err
	}

	return SettingsDetail, nil
}

func (menu *MenuModel) SettingsUpdates(settingsdetails TblGoTemplateSettings, DB *gorm.DB) (err error) {

	var settinglist TblGoTemplateSeo

	result := DB.Table("tbl_go_template_settings").Where("tenant_id = ? and website_id=? and template_id =?", settingsdetails.TenantId, settingsdetails.WebsiteId, settingsdetails.TemplateID).First(&settinglist)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {

			if err := DB.Table("tbl_go_template_settings").Create(&settingsdetails).Error; err != nil {

				return err
			}
		}

	} else {

		if settingsdetails.SiteName != "" {
			fmt.Println("settingsdetails.SiteName ", settingsdetails.TemplateID, settingsdetails.HeaderThame)
			if err := DB.Debug().Table("tbl_go_template_settings").Where("tenant_id = ? and website_id=? and template_id =?", settingsdetails.TenantId, settingsdetails.WebsiteId, settingsdetails.TemplateID).UpdateColumns(map[string]interface{}{"site_name": settingsdetails.SiteName, "site_logo": settingsdetails.SiteLogo, "site_logo_path": settingsdetails.SiteLogoPath, "site_fav_icon": settingsdetails.SiteFavIcon, "site_fav_icon_path": settingsdetails.SiteFavIconPath, "header_thame": settingsdetails.HeaderThame}).Error; err != nil {

				return err

			}

		}

		if len(settingsdetails.SocialMediaLink) != 0 {

			if err := DB.Table("tbl_go_template_settings").Where("tenant_id = ? and website_id=? and template_id =?", settingsdetails.TenantId, settingsdetails.WebsiteId, settingsdetails.TemplateID).UpdateColumns(map[string]interface{}{"social_media_link": settingsdetails.SocialMediaLink}).Error; err != nil {

				return err

			}

		}

		if settingsdetails.TemplateType != nil && settingsdetails.WebsiteUrl != "" {

			if err := DB.Table("tbl_go_template_settings").Where("tenant_id = ? and website_id=? and template_id =?", settingsdetails.TenantId, settingsdetails.WebsiteId, settingsdetails.TemplateID).UpdateColumns(map[string]interface{}{"template_type": settingsdetails.TemplateType, "website_url": settingsdetails.WebsiteUrl}).Error; err != nil {

				return err

			}
		}
	}

	return nil
}
