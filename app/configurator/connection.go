/*
User settings. Probably in Sqlite3? :thinking:
*/

package configurator

import (
	"darkbot/app/configurator/models"
	"darkbot/app/settings"
	"darkbot/app/settings/darkbot_logus"
	"darkbot/app/settings/types"

	"github.com/darklab8/darklab_goutils/goutils/logus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Configurator struct {
	db     *gorm.DB
	dbpath types.Dbpath
}

func (cg Configurator) GetClient() *gorm.DB {
	return cg.db
}
func (cg Configurator) GetDbpath() types.Dbpath {
	return cg.dbpath
}

func NewConfigurator(dbpath types.Dbpath) *Configurator {
	db, err := gorm.Open(
		sqlite.Open(string(dbpath)+"?cache=shared&mode=rwc&_journal_mode=WAL"), &gorm.Config{},
	)
	darkbot_logus.Log.CheckFatal(err, "failed to connect database at dbpath=", darkbot_logus.Dbpath(settings.Dbpath))

	return &Configurator{db: db, dbpath: dbpath}
}

func (cg *Configurator) AutoMigrateSchema() *Configurator {
	err := cg.db.AutoMigrate(
		&models.Channel{},
		&models.TagBase{},
		&models.TagPlayerFriend{},
		&models.TagPlayerEnemy{},
		&models.TagSystem{},
		&models.TagRegion{},
		&models.TagForumPostTrack{},
		&models.TagForumPostIgnore{},
		&models.TagForumSubforumTrack{},
		&models.TagForumSubforumIgnore{},
		&models.TagPlayerEvent{},
		&models.AlertNeutralPlayersEqualOrGreater{},
		&models.AlertEnemiesEqualOrGreater{},
		&models.AlertFriendsEqualOrGreater{},
		&models.AlertBaseHealthLowerThan{},
		&models.AlertBaseIfHealthDecreasing{},
		&models.AlertBaseIfUnderAttack{},
		&models.AlertPingMessage{},
		&models.ConfigBaseOrderingKey{},
	)
	if !darkbot_logus.Log.CheckWarn(err, "AutoMigrateSchema was executed with problems", logus.OptError(err)) {
		darkbot_logus.Log.Info("AutoMigrateSchema was executed fine")
	}
	return cg
}
