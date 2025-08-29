package menu

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TblGoTemplates struct {
	Id                 int
	TemplateName       string
	TemplateImage      string
	IsDeleted          int
	TenantId           string
	CreatedBy          int
	ChannelSlugName    string
	TemplateModuleName string
	CreatedOn          time.Time
	DateString         string `gorm:"-:migration;<-:false"`
}

func (menu *MenuModel) ListGoTemplates(tenantid string, DB *gorm.DB) (list []TblGoTemplates, count int64, err error) {
	var GoTemplatesList []TblGoTemplates
	query := DB.Table("tbl_go_templates").Where("is_deleted = 0")

	if tenantid != "" {
		query = query.Where("tenant_id = ?", tenantid)
	} else {
		query = query.Where("tenant_id IS NULL")
	}

	if err = query.Count(&count).Error; err != nil {
		return []TblGoTemplates{}, 0, err
	}

	if err = query.Order("id ASC").Find(&GoTemplatesList).Error; err != nil {
		return []TblGoTemplates{}, 0, err
	}

	return GoTemplatesList, count, nil
}

func (menu *MenuModel) GetTemplateById(moduleid int, tenantid string, DB *gorm.DB) (templates TblGoTemplates, err error) {

	if err := DB.Table("tbl_go_templates").Where("is_deleted = 0 and id=? ", moduleid).Order("id asc").Find(&templates).Error; err != nil {

		return TblGoTemplates{}, err
	}

	return templates, nil
}

func (menu *MenuModel) CloneTemplatesBySlug(db *gorm.DB, slug string, tenantID string, userid int, usertype string) error {
	var templates []TblGoTemplates

	switch usertype {

	case "new":

		if err := db.Debug().Where("channel_slug_name = ? AND is_deleted = 0 AND tenant_id IS NULL", slug).Find(&templates).Error; err != nil {
			return fmt.Errorf("error fetching templates for new user: %w", err)
		}

	case "old":

		var count int64
		if err := db.Model(&TblGoTemplates{}).Where("channel_slug_name = ? AND tenant_id = ? AND is_deleted = 0", slug, tenantID).Count(&count).Error; err != nil {
			return fmt.Errorf("error checking existing old user templates: %w", err)
		}

		if count > 0 {

			return nil
		}

		if err := db.Debug().Where("channel_slug_name = ? AND is_deleted = 0 AND tenant_id IS NULL", slug).Find(&templates).Error; err != nil {
			return fmt.Errorf("error fetching global templates for old user: %w", err)
		}

	default:
		return fmt.Errorf("invalid userType: must be 'new' or 'old'")
	}

	if len(templates) == 0 {
		return fmt.Errorf("no templates found to clone for slug: %s", slug)
	}

	var newTemplates []TblGoTemplates
	now := time.Now().UTC()

	for _, t := range templates {

		t.Id = 0
		t.TenantId = tenantID
		t.CreatedOn = now
		t.CreatedBy = userid
		newTemplates = append(newTemplates, t)

	}

	if err := db.Create(&newTemplates).Error; err != nil {
		return fmt.Errorf("error inserting templates: %w", err)
	}
	var template TblGoTemplates
	if err := db.Debug().Where("is_deleted = 0 AND tenant_id=? and created_by=?", tenantID, userid).First(&template).Error; err != nil {
		return fmt.Errorf("error fetching templates for new user: %w", err)
	}

	if err := db.Debug().Table("tbl_users").Where("is_deleted = 0 AND tenant_id = ? AND id = ?", tenantID, userid).Updates(map[string]interface{}{"go_template_default": template.Id}).Error; err != nil {
		return err
	}

	return nil
}
