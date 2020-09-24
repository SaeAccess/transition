package transition

import (
	"fmt"
	"strings"

	"github.com/qor/audited"
	"gorm.io/gorm"
)

// StateChangeLog a model that used to keep state change logs
type StateChangeLog struct {
	gorm.Model
	ReferTable string
	ReferID    string
	From       string
	To         string
	Note       string `sql:"size:1024"`
	audited.AuditedModel
}

// GenerateReferenceKey generate reference key used for change log
func GenerateReferenceKey(model interface{}, db *gorm.DB) string {
	var (
		primaryValues []string
	)

	for _, field := range db.PrimaryFields() {
		primaryValues = append(primaryValues, fmt.Sprint(field.Field.Interface()))
	}

	return strings.Join(primaryValues, "::")
}

// GetStateChangeLogs get state change logs
func GetStateChangeLogs(model interface{}, db *gorm.DB) []StateChangeLog {
	var (
		changelogs []StateChangeLog
	)

	db.Where("refer_table = ? AND refer_id = ?", db.TableName(), GenerateReferenceKey(model, db)).Find(&changelogs)

	return changelogs
}

// GetLastStateChange gets last state change
func GetLastStateChange(model interface{}, db *gorm.DB) *StateChangeLog {
	var (
		changelog StateChangeLog
	)

	db.Where("refer_table = ? AND refer_id = ?", db.TableName(), GenerateReferenceKey(model, db)).Last(&changelog)
	if changelog.To == "" {
		return nil
	}
	return &changelog
}
