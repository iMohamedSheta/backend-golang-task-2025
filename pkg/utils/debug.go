package utils

import (
	"gorm.io/gorm"
)

// returns the SQL statement and the values of the statement
func DebugSQL(db *gorm.DB, model interface{}) {
	dryRunDB := db.Session(&gorm.Session{DryRun: true})
	stmt := dryRunDB.Find(model).Statement

	Dump("SQL: " + stmt.SQL.String())
	Dump("Vars:", stmt.Vars)
}
