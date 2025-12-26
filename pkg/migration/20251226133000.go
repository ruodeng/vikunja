package migration

import (
	"time"

	"src.techknowlogick.com/xormigrate"
	"xorm.io/xorm"
)

type Projects20251226133000 struct {
	EndDate time.Time `xorm:"DATETIME INDEX null"`
}

func (*Projects20251226133000) TableName() string {
	return "projects"
}

func init() {
	migrations = append(migrations, &xormigrate.Migration{
		ID:          "20251226133000",
		Description: "Add end date to projects",
		Migrate: func(tx *xorm.Engine) error {
			return tx.Sync(&Projects20251226133000{})
		},
		Rollback: func(tx *xorm.Engine) error {
			return nil
		},
	})
}
