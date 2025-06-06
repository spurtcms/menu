package menu

import (
	"errors"
	"fmt"

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
}

func (menu *MenuModel) SettingDetail(tenantid string, DB *gorm.DB) (setting TblGoTemplateSettings, err error) {

	var SettingsDetail TblGoTemplateSettings

	if err := DB.Table("tbl_go_template_settings").Where("tenant_id = ?", tenantid).First(&SettingsDetail).Error; err != nil {

		return TblGoTemplateSettings{}, err
	}

	return SettingsDetail, nil
}

func (menu *MenuModel) SettingsUpdates(settingsdetails TblGoTemplateSettings, DB *gorm.DB) (err error) {

	var settinglist TblGoTemplateSeo

	fmt.Println("hello")

	result := DB.Table("tbl_go_template_settings").Where("tenant_id = ?", settingsdetails.TenantId).First(&settinglist)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {

			if err := DB.Table("tbl_go_template_settings").Create(&settingsdetails).Error; err != nil {

				return err
			}

			return nil

		} else {

			fmt.Println("Database error:", result.Error)

		}

	} else {

		if settingsdetails.SiteName != "" {

			if err := DB.Table("tbl_go_template_settings").Where("tenant_id = ?", settingsdetails.TenantId).UpdateColumns(map[string]interface{}{"site_name": settingsdetails.SiteName}).Error; err != nil {

				return err

			}

		} else if settingsdetails.SiteLogo != "" {

			if err := DB.Table("tbl_go_template_settings").Where("tenant_id = ?", settingsdetails.TenantId).UpdateColumns(map[string]interface{}{"site_logo": settingsdetails.SiteLogo, "site_logo_path": settingsdetails.SiteLogoPath}).Error; err != nil {

				return err
			}

		} else if settingsdetails.SiteFavIcon != "" {

			if err := DB.Table("tbl_go_template_settings").Where("tenant_id = ?", settingsdetails.TenantId).UpdateColumns(map[string]interface{}{"site_fav_icon": settingsdetails.SiteFavIcon, "site_fav_icon_path": settingsdetails.SiteFavIconPath}).Error; err != nil {

				return err
			}

		} else if settingsdetails.WebsiteUrl != "" {

			if err := DB.Table("tbl_go_template_settings").Where("tenant_id = ?", settingsdetails.TenantId).UpdateColumns(map[string]interface{}{"website_url": settingsdetails.WebsiteUrl}).Error; err != nil {

				return err
			}

			if err := DB.Table("tbl_users").Where("tenant_id = ?", settingsdetails.TenantId).UpdateColumns(map[string]interface{}{"subdomain": settingsdetails.WebsiteUrl}).Error; err != nil {

				return err
			}

		}

		return nil

	}

	return nil
}
