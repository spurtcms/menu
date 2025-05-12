package postgres

import (
	"time"

	"gorm.io/gorm"
)

type TblMenus struct {
	Id          int       `gorm:"primaryKey;auto_increment;type:serial"`
	Name        string    `gorm:"type:character varying"`
	Description string    `gorm:"type:character varying"`
	TenantId    string    `gorm:"type:character varying"`
	CreatedOn   time.Time `gorm:"type:timestamp without time zone"`
	CreatedBy   int       `gorm:"type:integer"`
	IsDeleted   int       `gorm:"type:integer"`
	DeletedOn   time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy   int       `gorm:"DEFAULT:NULL"`
	ModifiedOn  time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy  int       `gorm:"DEFAULT:NULL"`
	Status      int       ` gorm:"type:integer"`
	UrlPath     string    `gorm:"type:character varying"`
	ParentId    int       `gorm:"type:integer"`
}

func MigrationTables(db *gorm.DB) {

	db.AutoMigrate(
		&TblMenus{},
	)

}
