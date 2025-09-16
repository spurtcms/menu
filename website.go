package menu

import "time"

func (menu *Menu) CreateWebsite(websiteinfo TblWebsite) ( TblWebsite,error ){

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblWebsite{}, AuthError
	}

	websiteinfo.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	website,err := menumodel.CreateWebsite(&websiteinfo, menu.DB)

	if err != nil {

		return TblWebsite{},err
	}

	return website, nil
}

func (menu *Menu) WebsiteList(limit int, offset int, filter Filter, tenantid string) (websitelist []TblWebsite, count int64, err error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return []TblWebsite{}, 0, AuthError
	}

	menumodel.DataAccess = menu.DataAccess
	menumodel.Userid = menu.UserId

	_, totalcount, _ := menumodel.WebsiteList(0, 0, filter, menu.DB, tenantid)

	weblist, _, cerr := menumodel.WebsiteList(limit, offset, filter, menu.DB, tenantid)

	if cerr != nil {

		return []TblWebsite{}, 0, cerr
	}

	return weblist, totalcount, nil

}

/*UpdateWebsite*/
func (menu *Menu) UpdateWebsite(req TblWebsite) (TblWebsite, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblWebsite{}, AuthError
	}

	if req.Id <= 0 || req.Name == "" {

		return TblWebsite{}, ErrorMenuName
	}

	req.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	updatemenu, err := menumodel.UpdateWebsite(&req, menu.DB)

	if err != nil {

		return TblWebsite{}, err
	}

	return updatemenu, nil

}

func (menu *Menu) GetWebsiteById(id int, tenantid string) (TblWebsite, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblWebsite{}, AuthError
	}
	website, err := menumodel.GetWebsiteById(id, tenantid, menu.DB)

	if err != nil {

		return TblWebsite{}, err

	}
	return website, nil
}

func (menu *Menu) GetWebsiteByName(name string) (TblWebsite, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblWebsite{}, AuthError
	}
	website, err := menumodel.GetWebsiteByName(name, menu.DB)

	if err != nil {

		return TblWebsite{}, err

	}
	return website, nil
}

func (menu *Menu) CheckSiteName(sitename string,webid int) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}
	err := menumodel.CheckSiteName(sitename, webid, menu.DB)

	if err != nil {

		return err

	}

	return nil
}
