package menu

import "fmt"

func (menu *Menu) GoTemplatesList(tenantid string) (goTemplateList []TblGoTemplates, count int64, err error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return []TblGoTemplates{}, 0, AuthError
	}

	list, count, err := menumodel.ListGoTemplates(tenantid, menu.DB)

	if err != nil {

		return []TblGoTemplates{}, 0, err

	}
	fmt.Println("Hello")

	return list, count, nil
}

func (menu *Menu) GetTemplateById(id int, tenantid string) (TblGoTemplates, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblGoTemplates{}, AuthError
	}
	template, err := menumodel.GetTemplateById(id, tenantid, menu.DB)

	if err != nil {

		return TblGoTemplates{}, err

	}
	return template, nil
}

func (menu *Menu) CloneTemplatesBySlug(slug string, tenantid string, userid int,usertype string) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}
	err := menumodel.CloneTemplatesBySlug(menu.DB, slug, tenantid, userid,usertype)

	if err != nil {

		return err

	}
	return nil
}
