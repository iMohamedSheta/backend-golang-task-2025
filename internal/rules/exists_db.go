package rules

import (
	"strings"
	"taskgo/internal/deps"

	"github.com/go-playground/validator/v10"
)

// ExistsDB checks if a given column in a given table exists in the database example: "exists_db=table-column"
func ExistsDB(fl validator.FieldLevel) bool {
	db := deps.Gorm().DB
	log := deps.Log().Log()
	param := fl.Param() // expected format: "table-column"
	parts := strings.Split(param, "-")
	if len(parts) != 2 {
		log.Error("Invalid unique validator format, expected: table-column")
		return false
	}

	tableName := strings.TrimSpace(parts[0])
	columnName := strings.TrimSpace(parts[1])
	if tableName == "" || columnName == "" {
		return false
	}

	var count int64
	err := db.Table(tableName).Where(columnName+" = ?", fl.Field().String()).Count(&count).Error
	if err != nil {
		log.Error("Error while checking unique constraint: " + err.Error())
		return false
	}

	return count == 1
}
