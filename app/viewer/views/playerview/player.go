package playerview

import (
	_ "embed"
	"fmt"
	"text/template"

	"github.com/darklab8/fl-darkbot/app/scrappy/player"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/fl-darkbot/app/viewer/apis"
	"github.com/darklab8/fl-darkbot/app/viewer/views"
	"github.com/darklab8/fl-darkbot/app/viewer/views/viewer_msg"

	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils"
	"github.com/darklab8/go-utils/utils/utils_types"
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
	mainTable *views.ViewTable
	alertTmpl *views.ViewTable
}
type PlayersEnemies struct {
	mainTable *views.ViewTable
	alertTmpl *views.ViewTable
}
type PlayersNeutral struct {
	mainTable *views.ViewTable
	alertTmpl *views.ViewTable
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
	templator.channelID = channelID //
	templator.friends.mainTable = views.NewViewTable(viewer_msg.NewTableMsg(
		types.ViewID("#darkbot-players-friends-table"),
		types.ViewHeader("**Friend players in all systems and regions**\n"),
		types.ViewBeginning("```diff\n"),
		types.ViewEnd("```\n"),
	))
	templator.neutral.mainTable = views.NewViewTable(viewer_msg.NewTableMsg(
		types.ViewID("#darkbot-players-neutral-table"),
		types.ViewHeader("**Neutral players in tracked systems and regions**\n"),
		types.ViewBeginning("```json\n"),
		types.ViewEnd("```\n"),
	))
	templator.enemies.mainTable = views.NewViewTable(viewer_msg.NewTableMsg(
		types.ViewID("#darkbot-players-enemies-table"),
		types.ViewHeader("**Enemy players in tracked systems and regions**\n"),
		types.ViewBeginning("```diff\n"),
		types.ViewEnd("```\n"),
	))
	templator.friends.alertTmpl = views.NewViewTable(viewer_msg.NewAlertMsg(
		types.ViewID("#darkbot-players-friends-alert"),
	))
	templator.neutral.alertTmpl = views.NewViewTable(viewer_msg.NewAlertMsg(
		types.ViewID("#darkbot-players-neutral-alert"),
	))
	templator.enemies.alertTmpl = views.NewViewTable(viewer_msg.NewAlertMsg(
		types.ViewID("#darkbot-players-enemies-alert"),
	))

	templator.SharedViewTableSplitter = views.NewSharedViewSplitter(
		api,
		channelID,
		&templator,
		templator.friends.mainTable,
		templator.neutral.mainTable,
		templator.enemies.mainTable,
		templator.friends.alertTmpl,
		templator.neutral.alertTmpl,
		templator.enemies.alertTmpl,
	)
	return &templator
}

func (t *PlayersTemplates) GenerateRecords() error {
	record, err := t.api.Scrappy.GetPlayerStorage().GetLatestRecord()
	if logus.Log.CheckWarn(err, "unable to get players") {
		return err
	}

	systemTags, _ := t.api.Players.Systems.TagsList(t.channelID)
	regionTags, _ := t.api.Players.Regions.TagsList(t.channelID)
	friendTags, _ := t.api.Players.Friends.TagsList(t.channelID)
	enemyTags, _ := t.api.Players.Enemies.TagsList(t.channelID)
	logus.Log.Debug(
		"PlayerTemplatesRender next",
		typelog.Items("systemTags", systemTags),
		typelog.Items("friendTags", friendTags),
		typelog.Items("enemyTags", enemyTags),
		typelog.Items("record.List", record.List),
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

	logus.Log.Debug("friendPlayers=", typelog.Items("friendPlayers", friendPlayers))
	logus.Log.Debug("enemyPlayers=", typelog.Items("enemyPlayers", enemyPlayers))
	logus.Log.Debug("neutralPlayers=", typelog.Items("neutralPlayers", neutralPlayers))

	protectAgainstResend := func(player *[]player.Player, view *views.ViewTable) {
		if len(*player) == 0 {
			view.AppendRecord(types.ViewRecord(" \n"))
		}
	}

	if len(systemTags) > 0 || len(regionTags) > 0 {
		for _, playerVars := range neutralPlayers {
			t.neutral.mainTable.AppendRecord(types.ViewRecord(utils.TmpRender(playerTemplate, playerVars)))
		}

		protectAgainstResend(&neutralPlayers, t.neutral.mainTable)
	}

	if (len(systemTags) > 0 || len(regionTags) > 0) && len(enemyTags) > 0 {
		for _, playerVars := range enemyPlayers {
			t.enemies.mainTable.AppendRecord(types.ViewRecord(fmt.Sprintf("-%s", utils.TmpRender(playerTemplate, playerVars))))
		}

		protectAgainstResend(&enemyPlayers, t.enemies.mainTable)
	}

	if len(friendTags) > 0 {
		for _, playerVars := range friendPlayers {
			t.friends.mainTable.AppendRecord(types.ViewRecord(fmt.Sprintf("+%s", utils.TmpRender(playerTemplate, playerVars))))
		}

		protectAgainstResend(&friendPlayers, t.friends.mainTable)
	}

	// Alerts

	if alertNeutralCount, err := t.api.Alerts.NeutralsGreaterThan.Status(t.channelID); err == nil {
		if len(neutralPlayers) >= alertNeutralCount {
			t.neutral.alertTmpl.SetHeader(views.RenderAlertTemplate(t.channelID, fmt.Sprintf("Amount %d of neutral players is above or equal threshold %d", len(neutralPlayers), alertNeutralCount), t.api))
			t.neutral.alertTmpl.AppendRecord("")
		}
	}
	if alertEnemyCount, err := t.api.Alerts.EnemiesGreaterThan.Status(t.channelID); err == nil {
		if len(enemyPlayers) >= alertEnemyCount {
			t.enemies.alertTmpl.SetHeader(views.RenderAlertTemplate(t.channelID, fmt.Sprintf("Amount %d of enemy players is above or equal threshold %d", len(enemyPlayers), alertEnemyCount), t.api))
			t.enemies.alertTmpl.AppendRecord("")
		}
	}
	if alertFriendCount, err := t.api.Alerts.FriendsGreaterThan.Status(t.channelID); err == nil {
		if len(friendPlayers) >= alertFriendCount {
			t.friends.alertTmpl.SetHeader(views.RenderAlertTemplate(t.channelID, fmt.Sprintf("Amount %d of friendly players is above or equal threshold %d", len(friendPlayers), alertFriendCount), t.api))
			t.friends.alertTmpl.AppendRecord("")
		}
	}
	return nil
}

//go:embed player_template.md
var playerMarkup utils_types.TemplateExpression
var playerTemplate *template.Template

func init() {
	playerTemplate = utils.TmpInit(playerMarkup)
}
