package menu

import (
	"encoding/json"
	"fmt"
)

func (menu *Menu) SettingsDetail(tenantid string, websiteid int) (setting TblGoTemplateSettings, err error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblGoTemplateSettings{}, AuthError
	}

	settingsdetail, err := menumodel.SettingDetail(tenantid, websiteid, menu.DB)

	if err != nil {

		return TblGoTemplateSettings{}, err

	}

	if len(settingsdetail.SocialMediaLink) > 0 {
		_ = json.Unmarshal(settingsdetail.SocialMediaLink, &settingsdetail.SocialLinks)
	}

	return settingsdetail, nil
}

func (menu *Menu) SettingUpdate(settingsdetails TblGoTemplateSettings) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}

	Settings := TblGoTemplateSettings{

		TemplateID: settingsdetails.TemplateID,
		SiteName:        settingsdetails.SiteName,
		SiteLogo:        settingsdetails.SiteLogo,
		SiteLogoPath:    settingsdetails.SiteLogoPath,
		SiteFavIcon:     settingsdetails.SiteFavIcon,
		SiteFavIconPath: settingsdetails.SiteFavIconPath,
		WebsiteUrl:      settingsdetails.WebsiteUrl,
		TenantId:        settingsdetails.TenantId,
		WebsiteId:       settingsdetails.WebsiteId,
		TemplateType:    settingsdetails.TemplateType,
		SocialMediaLink: settingsdetails.SocialMediaLink,
		HeaderThame:     settingsdetails.HeaderThame,
	}

	fmt.Println("")

	err := menumodel.SettingsUpdates(Settings, menu.DB)

	if err != nil {

		return err

	}

	return nil
}
func (menu *Menu) SettingDetailBasedONTemp(TamplateID string, tenantid string, websiteid int) (setting TblGoTemplateSettings, err error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblGoTemplateSettings{}, AuthError
	}

	settingsdetail, err := menumodel.SettingDetailBasedONTemp(TamplateID, tenantid, websiteid, menu.DB)

	if err != nil {

		return TblGoTemplateSettings{}, err

	}

	if len(settingsdetail.SocialMediaLink) > 0 {
		_ = json.Unmarshal(settingsdetail.SocialMediaLink, &settingsdetail.SocialLinks)
	}

	return settingsdetail, nil

}
