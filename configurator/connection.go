/*
User settings. Probably in Sqlite3? :thinking:
*/

package configurator

import (
	"darkbot/configurator/models"
	"darkbot/settings"
	"darkbot/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Configurator struct {
	db *gorm.DB
}

func (cg Configurator) GetClient() *gorm.DB {
	return cg.db
}

func NewConfigurator() Configurator {
	db, err := gorm.Open(sqlite.Open(settings.Dbpath+"?cache=shared&mode=rwc&_journal_mode=WAL"), &gorm.Config{})
	utils.CheckPanic(err, "failed to connect database")

	// Auto Migrate the schema
	db.AutoMigrate(
		&models.Channel{},
		&models.TagBase{},
		&models.TagPlayerFriend{},
		&models.TagPlayerEnemy{},
		&models.TagSystem{},
		&models.TagRegion{},
		&models.TagForumPostTrack{},
		&models.TagForumPostIgnore{},
		&models.AlertPlayerUnrecognized{},
		&models.AlertPlayerEnemy{},
		&models.AlertPlayerFriend{},
		&models.AlertBase{},
	)

	return Configurator{db: db}
}
