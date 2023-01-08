/*
User settings. Probably in Sqlite3? :thinking:
*/

package configurator

import (
	"darkbot/settings"
	"darkbot/utils"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetConnection() *gorm.DB {
	workdir := filepath.Dir(utils.GetCurrrentFolder())
	dbpath := filepath.Join(workdir, "data", settings.Config.ConfiguratorDbname+".sqlite3")
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	utils.CheckPanic(err, "failed to connect database")

	// Auto Migrate the schema
	db.AutoMigrate(&Channel{})

	return db
}
