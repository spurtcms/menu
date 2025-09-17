package menu

import (
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

	return settingsdetail, nil
}

func (menu *Menu) SettingUpdate(settingsdetails TblGoTemplateSettings) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}

	Settings := TblGoTemplateSettings{
		SiteName:        settingsdetails.SiteName,
		SiteLogo:        settingsdetails.SiteLogo,
		SiteLogoPath:    settingsdetails.SiteLogoPath,
		SiteFavIcon:     settingsdetails.SiteFavIcon,
		SiteFavIconPath: settingsdetails.SiteFavIconPath,
		WebsiteUrl:      settingsdetails.WebsiteUrl,
		TenantId:        settingsdetails.TenantId,
		WebsiteId: settingsdetails.WebsiteId,
	}

	fmt.Println("hello::", Settings)

	err := menumodel.SettingsUpdates(Settings, menu.DB)

	if err != nil {

		return err

	}

	return nil
}
