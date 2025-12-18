package migration

import (
	"fmt"

	"src.techknowlogick.com/xormigrate"
	"xorm.io/xorm"
)

type ShiftProjectRights20251217150842 struct{}

func (ShiftProjectRights20251217150842) TableName() string {
	return "project_users"
}

func init() {
	migrations = append(migrations, &xormigrate.Migration{
		ID:          "20251217150842",
		Description: "Shift project rights to introduce Restricted(1) tier",
		Migrate: func(tx *xorm.Engine) error {
			return performRightShift(tx, true)
		},
		Rollback: func(tx *xorm.Engine) error {
			return performRightShift(tx, false)
		},
	})
}

func performRightShift(tx *xorm.Engine, isUp bool) error {
	targetTables := []string{
		"project_users",
		"team_projects",
		"link_shares",
	}

	for _, tableName := range targetTables {
		exists, err := tx.IsTableExist(tableName)
		if err != nil {
			return fmt.Errorf("failed to check existence of table %s: %w", tableName, err)
		}
		if !exists {
			continue
		}

		driverName := tx.DriverName()
		quoteChar := "`"
		if driverName == "postgres" {
			quoteChar = "\""
		}

		quotedTable := fmt.Sprintf("%s%s%s", quoteChar, tableName, quoteChar)
		quotedColumn := fmt.Sprintf("%spermission%s", quoteChar, quoteChar)

		var query string
		if isUp {
			// Upgrade: >= 1 becomes +1
			// 1 (Write) -> 2 (Write)
			// 2 (Admin) -> 3 (Admin)
			// 0 (Read) -> 0 (Read)
			query = fmt.Sprintf(
				"UPDATE %s SET %s = %s + 1 WHERE %s >= 1",
				quotedTable, quotedColumn, quotedColumn, quotedColumn,
			)
		} else {
			// Rollback:
			// 1. Reset any newly created Restricted(1) to Read(0) for safety
			safetyQuery := fmt.Sprintf(
				"UPDATE %s SET %s = 0 WHERE %s = 1",
				quotedTable, quotedColumn, quotedColumn,
			)
			if _, err := tx.Exec(safetyQuery); err != nil {
				return fmt.Errorf("rollback safety step failed for table %s: %w", tableName, err)
			}

			// 2. Shift back >= 2 to -1
			query = fmt.Sprintf(
				"UPDATE %s SET %s = %s - 1 WHERE %s >= 2",
				quotedTable, quotedColumn, quotedColumn, quotedColumn,
			)
		}

		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("failed to migrate table %s: %w", tableName, err)
		}
	}

	return nil
}
