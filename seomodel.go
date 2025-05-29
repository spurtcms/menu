package menu

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TblGoTemplateSeo struct {
	Id               int
	PageTitle        string
	PageDescription  string
	PageKeyword      string
	StoreTitle       string
	StoreDescription string
	StoreKeyword     string
	SiteMapName      string
	SiteMapPath      string
	TenantId         string
}

func (menu *MenuModel) SeoDetails(tenantid string, DB *gorm.DB) (seo TblGoTemplateSeo, err error) {

	var SeoDetail TblGoTemplateSeo

	if err := DB.Table("tbl_go_template_seos").Where("tenant_id = ?", tenantid).First(&SeoDetail).Error; err != nil {

		return TblGoTemplateSeo{}, err
	}

	fmt.Println("Hello")

	return SeoDetail, nil
}

func (menu *MenuModel) SeoUpdates(seodetails TblGoTemplateSeo, DB *gorm.DB) (err error) {

	var seolist TblGoTemplateSeo

	result := DB.Table("tbl_go_template_seos").Where("tenant_id = ?", seodetails.TenantId).First(&seolist)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {

			if err := DB.Table("tbl_go_template_seos").Create(&seodetails).Error; err != nil {

				return err
			}

			return nil

		} else {

			fmt.Println("Database error:", result.Error)

		}

	} else {

		if seodetails.PageTitle != "" {

			if err := DB.Table("tbl_go_template_seos").Where("tenant_id = ?", seodetails.TenantId).UpdateColumns(map[string]interface{}{"page_title": seodetails.PageTitle, "page_description": seodetails.PageDescription, "page_keyword": seodetails.PageKeyword}).Error; err != nil {

				return err

			}

		} else if seodetails.StoreTitle != "" {

			if err := DB.Table("tbl_go_template_seos").Where("tenant_id = ?", seodetails.TenantId).UpdateColumns(map[string]interface{}{"store_title": seodetails.StoreTitle, "store_description": seodetails.StoreDescription, "store_keyword": seodetails.StoreKeyword}).Error; err != nil {

				return err
			}

		} else if seodetails.SiteMapName != "" {

			if err := DB.Table("tbl_go_template_seos").Where("tenant_id = ?", seodetails.TenantId).UpdateColumns(map[string]interface{}{"site_map_name": seodetails.SiteMapName, "site_map_path": seodetails.SiteMapPath}).Error; err != nil {

				return err
			}

		}

		return nil

	}

	return nil
}
