package menu

import (
	"fmt"

	"gorm.io/gorm"
)

type TblGoTemplates struct {
	Id            int
	TemplateName  string
	TemplateImage string
	IsDeleted     int
	TenantId      string
}

func (menu *MenuModel) ListGoTemplates(isdeleted int, DB *gorm.DB) (list []TblGoTemplates, count int64, err error) {

	var GoTemplatesList []TblGoTemplates

	if err := DB.Table("tbl_go_templates").Where("is_deleted=?", isdeleted).Find(&GoTemplatesList).Count(&count).Error; err != nil {

		return []TblGoTemplates{}, 0, err
	}

	fmt.Println("Hello")

	return GoTemplatesList, count, nil
}

func (menu *MenuModel) GetTemplateById(moduleid int, tenantid string, DB *gorm.DB) (templates TblGoTemplates, err error) {

	if err := DB.Debug().Table("tbl_go_templates").Where("is_deleted = 0 and id=? ", moduleid).Order("id asc").Find(&templates).Error; err != nil {

		return TblGoTemplates{}, err
	}

	return templates, nil
}
