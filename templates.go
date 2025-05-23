package menu

import "fmt"

func (menu *Menu) GoTemplatesList(isdeleted int) (goTemplateList []TblGoTemplates, count int64, err error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return []TblGoTemplates{}, 0, AuthError
	}

	list, count, err := menumodel.ListGoTemplates(isdeleted, menu.DB)

	if err != nil {

		return []TblGoTemplates{}, 0, err

	}
	fmt.Println("Hello")

	return list, count, nil
}
