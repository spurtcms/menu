package menu

import (
	"fmt"
)

func (menu *Menu) SeoDetail(tenantid string) (seo TblGoTemplateSeo, err error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblGoTemplateSeo{}, AuthError
	}

	seodetail, err := menumodel.SeoDetails(tenantid, menu.DB)

	if err != nil {

		return TblGoTemplateSeo{}, err

	}
	fmt.Println("Hello")

	return seodetail, nil
}

func (menu *Menu) SeoUpdate(seodetails TblGoTemplateSeo) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}

	SEO := TblGoTemplateSeo{
		PageTitle:        seodetails.PageTitle,
		PageDescription:  seodetails.PageDescription,
		PageKeyword:      seodetails.PageKeyword,
		StoreTitle:       seodetails.StoreTitle,
		StoreDescription: seodetails.StoreDescription,
		StoreKeyword:     seodetails.StoreKeyword,
		SiteMapName:      seodetails.SiteMapName,
		SiteMapPath:      seodetails.SiteMapPath,
		TenantId:         seodetails.TenantId,
	}

	fmt.Println("hello::", SEO)

	err := menumodel.SeoUpdates(SEO, menu.DB)

	if err != nil {

		return err

	}

	return nil
}
