package mysql

import (
	"time"

	"gorm.io/gorm"
)

type TblMenus struct {
	Id          int       `gorm:"primaryKey;auto_increment"`
	Name        string    `gorm:"type:varchar(255)"`
	Description string    `gorm:"type:varchar(255)"`
	TenantId    string    `gorm:"type:varchar(255)"`
	CreatedOn   time.Time `gorm:"type:datetime"`
	CreatedBy   int       `gorm:"type:integer"`
	IsDeleted   int       `gorm:"type:integer"`
	DeletedOn   time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	DeletedBy   int       `gorm:"type:integer;DEFAULT:NULL"`
	ModifiedOn  time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	ModifiedBy  int       `gorm:"type:integer;DEFAULT:NULL"`
	Status      int       ` gorm:"type:integer"`
	UrlPath     string    `gorm:"type:varchar(255)"`
	ParentId    int       `gorm:"type:integer"`
}

func MigrationTables(db *gorm.DB) {

	db.AutoMigrate(
		&TblMenus{},
	)

}
