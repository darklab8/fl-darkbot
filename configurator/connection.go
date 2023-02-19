/*
User settings. Probably in Sqlite3? :thinking:
*/

package configurator

import (
	"darkbot/configurator/models"
	"darkbot/dtypes"
	"darkbot/settings"
	"darkbot/utils/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Configurator struct {
	db *gorm.DB
}

func (cg Configurator) GetClient() *gorm.DB {
	return cg.db
}

func NewConfigurator(dbpath dtypes.Dbpath) Configurator {
	db, err := gorm.Open(
		sqlite.Open(string(dbpath)+"?cache=shared&mode=rwc&_journal_mode=WAL"), &gorm.Config{},
	)
	logger.CheckPanic(err, "failed to connect database at dbpath=", string(settings.Dbpath))

	return Configurator{db: db}
}

func (cg Configurator) Migrate() Configurator {
	// Auto Migrate the schema
	cg.db.AutoMigrate(
		&models.Channel{},
		&models.TagBase{},
		&models.TagPlayerFriend{},
		&models.TagPlayerEnemy{},
		&models.TagSystem{},
		&models.TagRegion{},
		&models.TagForumPostTrack{},
		&models.TagForumPostIgnore{},
		&models.AlertNeutralPlayersEqualOrGreater{},
		&models.AlertEnemiesEqualOrGreater{},
		&models.AlertFriendsEqualOrGreater{},
		&models.AlertBaseHealthLowerThan{},
		&models.AlertBaseIfHealthDecreasing{},
		&models.AlertBaseIfUnderAttack{},
		&models.AlertPingMessage{},
	)
	return cg
}
