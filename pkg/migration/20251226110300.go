package migration

import (
	"time"

	"src.techknowlogick.com/xormigrate"
	"xorm.io/xorm"
)

type Projects20251226110300 struct {
	ResponsibleUserID int64     `xorm:"bigint INDEX null"`
	EndDate           time.Time `xorm:"DATETIME INDEX null"`
}

func (Projects20251226110300) TableName() string {
	return "projects"
}

func init() {
	migrations = append(migrations, &xormigrate.Migration{
		ID:          "20251226110300",
		Description: "Add responsible user and end date to projects",
		Migrate: func(tx *xorm.Engine) error {
			return tx.Sync(Projects20251226110300{})
		},
		Rollback: func(tx *xorm.Engine) error {
			return nil
		},
	})
}
