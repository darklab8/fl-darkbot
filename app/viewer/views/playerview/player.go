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
	"strings"
	"text/template"
	"time"
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
	mainTable views.TemplateShared
	alertTmpl views.TemplateShared
}
type PlayersEnemies struct {
	mainTable views.TemplateShared
	alertTmpl views.TemplateShared
}
type PlayersNeutral struct {
	mainTable views.TemplateShared
	alertTmpl views.TemplateShared
}

type PlayersTemplates struct {
	friends PlayersFriends
	neutral PlayersNeutral
	enemies PlayersEnemies
	api     *apis.API
}

func NewTemplatePlayers(api *apis.API) *PlayersTemplates {
	templator := PlayersTemplates{}
	templator.api = api
	templator.friends.mainTable.Header = "#darkbot-players-friends-table"
	templator.neutral.mainTable.Header = "#darkbot-players-neutral-table"
	templator.enemies.mainTable.Header = "#darkbot-players-enemies-table"
	templator.friends.alertTmpl.Header = "#darkbot-players-friends-alert"
	templator.neutral.alertTmpl.Header = "#darkbot-players-neutral-alert"
	templator.enemies.alertTmpl.Header = "#darkbot-players-enemies-alert"
	return &templator
}

func TagContains(name string, tags []types.Tag) bool {
	for _, tag := range tags {
		if strings.Contains(name, string(tag)) {
			return true
		}
	}
	return false
}

func (t *PlayersTemplates) Render() error {
	record, err := t.api.Scrappy.GetPlayerStorage().GetLatestRecord()
	if logus.CheckWarn(err, "unable to get player msgs") {
		return err
	}

	systemTags, _ := t.api.Players.Systems.TagsList(t.api.ChannelID)
	regionTags, _ := t.api.Players.Regions.TagsList(t.api.ChannelID)
	friendTags, _ := t.api.Players.Friends.TagsList(t.api.ChannelID)
	enemyTags, _ := t.api.Players.Enemies.TagsList(t.api.ChannelID)
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
		if TagContains(player.Name, friendTags) {
			friendPlayers = append(friendPlayers, player)
			continue
		}

		if !TagContains(player.System, systemTags) && !TagContains(player.Region, regionTags) {
			continue
		}

		if TagContains(player.Name, enemyTags) {
			enemyPlayers = append(enemyPlayers, player)
			continue
		}

		neutralPlayers = append(neutralPlayers, player)
	}

	logus.Debug("friendPlayers=", logus.Items(friendPlayers, "friendPlayers"))
	logus.Debug("enemyPlayers=", logus.Items(enemyPlayers, "enemyPlayers"))
	logus.Debug("neutralPlayers=", logus.Items(neutralPlayers, "neutralPlayers"))

	if len(systemTags) > 0 || len(regionTags) > 0 {
		t.neutral.mainTable.Content = utils.TmpRender(playerTemplate, TemplateRendrerPlayerInput{
			Header:      t.neutral.mainTable.Header,
			LastUpdated: time.Now().String(),
			Players:     neutralPlayers,
			TableName:   "**Neutral players in tracked systems and regions**",
		})
	}

	if (len(systemTags) > 0 || len(regionTags) > 0) && len(enemyTags) > 0 {
		t.enemies.mainTable.Content = utils.TmpRender(playerTemplate, TemplateRendrerPlayerInput{
			Header:      t.enemies.mainTable.Header,
			LastUpdated: time.Now().String(),
			Players:     enemyPlayers,
			TableName:   "**Enemy players in tracked systems and regions**",
		})
	}

	if len(friendTags) > 0 {
		t.friends.mainTable.Content = utils.TmpRender(playerTemplate, TemplateRendrerPlayerInput{
			Header:      t.friends.mainTable.Header,
			LastUpdated: time.Now().String(),
			Players:     friendPlayers,
			TableName:   "**Friend players in all systems and regions**",
		})
	}

	// Alerts

	if alertNeutralCount, err := t.api.Alerts.NeutralsGreaterThan.Status(t.api.ChannelID); err == nil {
		if len(neutralPlayers) >= alertNeutralCount {
			t.neutral.alertTmpl.Content = views.RenderAlertTemplate(t.neutral.alertTmpl.Header, t.api.ChannelID, fmt.Sprintf("Amount %d of neutral players is above threshold %d", len(neutralPlayers), alertNeutralCount), t.api)
		}
	}
	if alertEnemyCount, err := t.api.Alerts.EnemiesGreaterThan.Status(t.api.ChannelID); err == nil {
		if len(enemyPlayers) >= alertEnemyCount {
			t.enemies.alertTmpl.Content = views.RenderAlertTemplate(t.enemies.alertTmpl.Header, t.api.ChannelID, fmt.Sprintf("Amount %d of enemy players is above threshold %d", len(enemyPlayers), alertEnemyCount), t.api)
		}
	}
	if alertFriendCount, err := t.api.Alerts.FriendsGreaterThan.Status(t.api.ChannelID); err == nil {
		if len(friendPlayers) >= alertFriendCount {
			t.friends.alertTmpl.Content = views.RenderAlertTemplate(t.friends.alertTmpl.Header, t.api.ChannelID, fmt.Sprintf("Amount %d of friendly players is above threshold %d", len(friendPlayers), alertFriendCount), t.api)
		}
	}
	return nil
}

func (t *PlayersTemplates) Send() {
	t.friends.mainTable.Send(t.api)
	t.neutral.mainTable.Send(t.api)
	t.enemies.mainTable.Send(t.api)

	t.friends.alertTmpl.Send(t.api)
	t.neutral.alertTmpl.Send(t.api)
	t.enemies.alertTmpl.Send(t.api)

}

func (t *PlayersTemplates) MatchMessageID(messageID types.DiscordMessageID) bool {

	if messageID == t.friends.mainTable.MessageID {
		return true
	}
	if messageID == t.neutral.mainTable.MessageID {
		return true
	}
	if messageID == t.enemies.mainTable.MessageID {
		return true
	}

	if messageID == t.friends.alertTmpl.MessageID {
		return true
	}
	if messageID == t.neutral.alertTmpl.MessageID {
		return true
	}
	if messageID == t.enemies.alertTmpl.MessageID {
		return true
	}

	return false
}

func (t *PlayersTemplates) DiscoverMessageID(content string, msgID types.DiscordMessageID) {
	if strings.Contains(content, t.friends.mainTable.Header) {
		t.friends.mainTable.MessageID = msgID
	}
	if strings.Contains(content, t.neutral.mainTable.Header) {
		t.neutral.mainTable.MessageID = msgID
	}
	if strings.Contains(content, t.enemies.mainTable.Header) {
		t.enemies.mainTable.MessageID = msgID
	}

	if strings.Contains(content, t.friends.alertTmpl.Header) {
		t.friends.alertTmpl.MessageID = msgID
	}
	if strings.Contains(content, t.neutral.alertTmpl.Header) {
		t.neutral.alertTmpl.MessageID = msgID
	}
	if strings.Contains(content, t.enemies.alertTmpl.Header) {
		t.enemies.alertTmpl.MessageID = msgID
	}
}

//go:embed player_template.md
var playerMarkup string
var playerTemplate *template.Template

func init() {
	playerTemplate = utils.TmpInit(playerMarkup)
}

type TemplateRendrerPlayerInput struct {
	Header      string
	LastUpdated string
	Players     []player.Player
	TableName   string
}
