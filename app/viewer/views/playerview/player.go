package playerview

import (
	"darkbot/app/scrappy/player"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
	"darkbot/app/viewer/apis"
	"darkbot/app/viewer/views"
	_ "embed"
	"fmt"
	"text/template"
)

// Discovery players-all, players-friends, players-enemies messages
// It can be one Struct :thinking:
// Query Configurator players
// Match friends across all map
// Match players in selected region and systems
// Match enemies in selected region and system and exclude them into separate set
// Render
// Send

type PlayersFriends struct {
	mainTable views.ViewTable
	alertTmpl views.ViewTable
}
type PlayersEnemies struct {
	mainTable views.ViewTable
	alertTmpl views.ViewTable
}
type PlayersNeutral struct {
	mainTable views.ViewTable
	alertTmpl views.ViewTable
}

type PlayersTemplates struct {
	friends PlayersFriends
	neutral PlayersNeutral
	enemies PlayersEnemies
	api     *apis.API
	*views.SharedViewTableSplitter
	channelID types.DiscordChannelID
}

func NewTemplatePlayers(api *apis.API, channelID types.DiscordChannelID) *PlayersTemplates {
	templator := PlayersTemplates{}
	templator.api = api
	templator.channelID = channelID
	templator.friends.mainTable.ViewID = "#darkbot-players-friends-table"
	templator.neutral.mainTable.ViewID = "#darkbot-players-neutral-table"
	templator.enemies.mainTable.ViewID = "#darkbot-players-enemies-table"
	templator.friends.alertTmpl.ViewID = "#darkbot-players-friends-alert"
	templator.neutral.alertTmpl.ViewID = "#darkbot-players-neutral-alert"
	templator.enemies.alertTmpl.ViewID = "#darkbot-players-enemies-alert"

	templator.SharedViewTableSplitter = views.NewSharedViewSplitter(
		api,
		channelID,
		&templator,
		&templator.friends.mainTable,
		&templator.neutral.mainTable,
		&templator.enemies.mainTable,
		&templator.friends.alertTmpl,
		&templator.neutral.alertTmpl,
		&templator.enemies.alertTmpl,
	)
	return &templator
}

func (t *PlayersTemplates) GenerateRecords() error {
	record, err := t.api.Scrappy.GetPlayerStorage().GetLatestRecord()
	if logus.CheckWarn(err, "unable to get players") {
		return err
	}

	systemTags, _ := t.api.Players.Systems.TagsList(t.channelID)
	regionTags, _ := t.api.Players.Regions.TagsList(t.channelID)
	friendTags, _ := t.api.Players.Friends.TagsList(t.channelID)
	enemyTags, _ := t.api.Players.Enemies.TagsList(t.channelID)
	logus.Debug(
		"PlayerTemplatesRender next",
		logus.Items(systemTags, "systemTags"),
		logus.Items(friendTags, "friendTags"),
		logus.Items(enemyTags, "enemyTags"),
		logus.Items(record.List, "record.List"),
	)
	neutralPlayers := []player.Player{}
	enemyPlayers := []player.Player{}
	friendPlayers := []player.Player{}

	for _, player := range record.List {
		if views.TagContains(player.Name, friendTags) {
			friendPlayers = append(friendPlayers, player)
			continue
		}

		if !views.TagContains(player.System, systemTags) && !views.TagContains(player.Region, regionTags) {
			continue
		}

		if views.TagContains(player.Name, enemyTags) {
			enemyPlayers = append(enemyPlayers, player)
			continue
		}

		neutralPlayers = append(neutralPlayers, player)
	}

	logus.Debug("friendPlayers=", logus.Items(friendPlayers, "friendPlayers"))
	logus.Debug("enemyPlayers=", logus.Items(enemyPlayers, "enemyPlayers"))
	logus.Debug("neutralPlayers=", logus.Items(neutralPlayers, "neutralPlayers"))

	protectAgainstResend := func(player *[]player.Player, view *views.ViewTable) {
		if len(*player) == 0 {
			view.AppendRecord(types.ViewRecord(" "))
		}
	}

	if len(systemTags) > 0 || len(regionTags) > 0 {
		t.neutral.mainTable.ViewBeginning = "**Neutral players in tracked systems and regions**\n```json\n"
		t.neutral.mainTable.ViewEnd = "```\n"
		for _, playerVars := range neutralPlayers {
			t.neutral.mainTable.AppendRecord(types.ViewRecord(utils.TmpRender(playerTemplate, playerVars)))
		}

		protectAgainstResend(&neutralPlayers, &t.neutral.mainTable)
	}

	if (len(systemTags) > 0 || len(regionTags) > 0) && len(enemyTags) > 0 {
		t.enemies.mainTable.ViewBeginning = "**Enemy players in tracked systems and regions**\n```diff\n"
		t.enemies.mainTable.ViewEnd = "```\n"

		for _, playerVars := range enemyPlayers {
			t.enemies.mainTable.AppendRecord(types.ViewRecord(fmt.Sprintf("-%s", utils.TmpRender(playerTemplate, playerVars))))
		}

		protectAgainstResend(&enemyPlayers, &t.enemies.mainTable)
	}

	if len(friendTags) > 0 {
		t.friends.mainTable.ViewBeginning = "**Friend players in all systems and regions**\n```diff\n"
		t.friends.mainTable.ViewEnd = "```\n"

		for _, playerVars := range friendPlayers {
			t.friends.mainTable.AppendRecord(types.ViewRecord(fmt.Sprintf("+%s", utils.TmpRender(playerTemplate, playerVars))))
		}

		protectAgainstResend(&friendPlayers, &t.friends.mainTable)
	}

	// Alerts

	if alertNeutralCount, err := t.api.Alerts.NeutralsGreaterThan.Status(t.channelID); err == nil {
		if len(neutralPlayers) >= alertNeutralCount {

			t.neutral.alertTmpl.AppendRecord(views.RenderAlertTemplate(
				t.channelID,
				fmt.Sprintf("Amount %d of neutral players is above threshold %d", len(neutralPlayers), alertNeutralCount),
				t.api,
			))
		}
	}
	if alertEnemyCount, err := t.api.Alerts.EnemiesGreaterThan.Status(t.channelID); err == nil {
		if len(enemyPlayers) >= alertEnemyCount {
			t.enemies.alertTmpl.AppendRecord(views.RenderAlertTemplate(
				t.channelID,
				fmt.Sprintf("Amount %d of enemy players is above threshold %d", len(enemyPlayers), alertEnemyCount),
				t.api,
			))
		}
	}
	if alertFriendCount, err := t.api.Alerts.FriendsGreaterThan.Status(t.channelID); err == nil {
		if len(friendPlayers) >= alertFriendCount {
			t.friends.alertTmpl.AppendRecord(views.RenderAlertTemplate(
				t.channelID,
				fmt.Sprintf("Amount %d of friendly players is above threshold %d", len(friendPlayers), alertFriendCount),
				t.api,
			))
		}
	}
	return nil
}

//go:embed player_template.md
var playerMarkup string
var playerTemplate *template.Template

func init() {
	playerTemplate = utils.TmpInit(playerMarkup)
}
