package models

import "github.com/goravel/framework/database/orm"

// AppSetting stores simple application-level key/value configuration.
// Keep this model intentionally small: it validates the ORM/SQLite integration
// without introducing a domain the template does not yet need.
type AppSetting struct {
	orm.Model
	Key   string
	Value string
}
