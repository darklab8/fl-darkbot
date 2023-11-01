package templ

import (
	"darkbot/app/scrappy/player"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
	"darkbot/app/viewer/apis"
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
	MainTable TemplateShared
	AlertTmpl TemplateShared
}
type PlayersEnemies struct {
	MainTable TemplateShared
	AlertTmpl TemplateShared
}
type PlayersNeutral struct {
	MainTable TemplateShared
	AlertTmpl TemplateShared
}

type PlayersTemplates struct {
	friends PlayersFriends
	neutral PlayersNeutral
	enemies PlayersEnemies
	API     apis.API
}

func NewTemplatePlayers(channelID types.DiscordChannelID, dbpath types.Dbpath) PlayersTemplates {
	templator := PlayersTemplates{}
	templator.API = apis.NewAPI(channelID, dbpath)
	templator.friends.MainTable.Header = "#darkbot-players-friends-table"
	templator.neutral.MainTable.Header = "#darkbot-players-neutral-table"
	templator.enemies.MainTable.Header = "#darkbot-players-enemies-table"
	templator.friends.AlertTmpl.Header = "#darkbot-players-friends-alert"
	templator.neutral.AlertTmpl.Header = "#darkbot-players-neutral-alert"
	templator.enemies.AlertTmpl.Header = "#darkbot-players-enemies-alert"
	return templator
}

func (b *PlayersTemplates) Setup(channelID types.DiscordChannelID) {
	b.API.ChannelID = channelID
	b.neutral.MainTable.MessageID = ""
	b.enemies.MainTable.MessageID = ""
	b.friends.MainTable.MessageID = ""
	b.neutral.MainTable.Content = ""
	b.enemies.MainTable.Content = ""
	b.friends.MainTable.Content = ""
	b.neutral.AlertTmpl.MessageID = ""
	b.enemies.AlertTmpl.MessageID = ""
	b.friends.AlertTmpl.MessageID = ""
	b.neutral.AlertTmpl.Content = ""
	b.enemies.AlertTmpl.Content = ""
	b.friends.AlertTmpl.Content = ""
}

func TagContains(name string, tags []string) bool {
	for _, tag := range tags {
		if strings.Contains(name, tag) {
			return true
		}
	}
	return false
}

func (t *PlayersTemplates) Render() {
	record, err := t.API.Scrappy.GetPlayerStorage().GetLatestRecord()
	if logus.CheckWarn(err, "unable to get player msgs") {
		return
	}

	systemTags, _ := t.API.Players.Systems.TagsList(t.API.ChannelID)
	regionTags, _ := t.API.Players.Regions.TagsList(t.API.ChannelID)
	friendTags, _ := t.API.Players.Friends.TagsList(t.API.ChannelID)
	enemyTags, _ := t.API.Players.Enemies.TagsList(t.API.ChannelID)
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
		t.neutral.MainTable.Content = utils.TmpRender(playerTemplate, TemplateRendrerPlayerInput{
			Header:      t.neutral.MainTable.Header,
			LastUpdated: time.Now().String(),
			Players:     neutralPlayers,
			TableName:   "**Neutral players in tracked systems and regions**",
		})
	}

	if (len(systemTags) > 0 || len(regionTags) > 0) && len(enemyTags) > 0 {
		t.enemies.MainTable.Content = utils.TmpRender(playerTemplate, TemplateRendrerPlayerInput{
			Header:      t.enemies.MainTable.Header,
			LastUpdated: time.Now().String(),
			Players:     enemyPlayers,
			TableName:   "**Enemy players in tracked systems and regions**",
		})
	}

	if len(friendTags) > 0 {
		t.friends.MainTable.Content = utils.TmpRender(playerTemplate, TemplateRendrerPlayerInput{
			Header:      t.friends.MainTable.Header,
			LastUpdated: time.Now().String(),
			Players:     friendPlayers,
			TableName:   "**Friend players in all systems and regions**",
		})
	}

	// Alerts

	if alertNeutralCount, _ := t.API.Alerts.NeutralsGreaterThan.Status(t.API.ChannelID); alertNeutralCount != nil {
		if len(neutralPlayers) >= *alertNeutralCount {
			t.neutral.AlertTmpl.Content = RenderAlertTemplate(t.neutral.AlertTmpl.Header, t.API.ChannelID, fmt.Sprintf("Amount %d of neutral players is above threshold %d", len(neutralPlayers), *alertNeutralCount), t.API)
		}
	}
	if alertEnemyCount, _ := t.API.Alerts.EnemiesGreaterThan.Status(t.API.ChannelID); alertEnemyCount != nil {
		if len(enemyPlayers) >= *alertEnemyCount {
			t.enemies.AlertTmpl.Content = RenderAlertTemplate(t.enemies.AlertTmpl.Header, t.API.ChannelID, fmt.Sprintf("Amount %d of enemy players is above threshold %d", len(enemyPlayers), *alertEnemyCount), t.API)
		}
	}
	if alertFriendCount, _ := t.API.Alerts.FriendsGreaterThan.Status(t.API.ChannelID); alertFriendCount != nil {
		if len(friendPlayers) >= *alertFriendCount {
			t.friends.AlertTmpl.Content = RenderAlertTemplate(t.friends.AlertTmpl.Header, t.API.ChannelID, fmt.Sprintf("Amount %d of friendly players is above threshold %d", len(friendPlayers), *alertFriendCount), t.API)
		}
	}
}

func (t *PlayersTemplates) Send() {
	t.friends.MainTable.Send(t.API)
	t.neutral.MainTable.Send(t.API)
	t.enemies.MainTable.Send(t.API)

	t.friends.AlertTmpl.Send(t.API)
	t.neutral.AlertTmpl.Send(t.API)
	t.enemies.AlertTmpl.Send(t.API)

}

func (t *PlayersTemplates) MatchMessageID(messageID types.DiscordMessageID) bool {

	if messageID == t.friends.MainTable.MessageID {
		return true
	}
	if messageID == t.neutral.MainTable.MessageID {
		return true
	}
	if messageID == t.enemies.MainTable.MessageID {
		return true
	}

	if messageID == t.friends.AlertTmpl.MessageID {
		return true
	}
	if messageID == t.neutral.AlertTmpl.MessageID {
		return true
	}
	if messageID == t.enemies.AlertTmpl.MessageID {
		return true
	}

	return false
}

func (t *PlayersTemplates) DiscoverMessageID(content string, msgID types.DiscordMessageID) {
	if strings.Contains(content, t.friends.MainTable.Header) {
		t.friends.MainTable.MessageID = msgID
	}
	if strings.Contains(content, t.neutral.MainTable.Header) {
		t.neutral.MainTable.MessageID = msgID
	}
	if strings.Contains(content, t.enemies.MainTable.Header) {
		t.enemies.MainTable.MessageID = msgID
	}

	if strings.Contains(content, t.friends.AlertTmpl.Header) {
		t.friends.AlertTmpl.MessageID = msgID
	}
	if strings.Contains(content, t.neutral.AlertTmpl.Header) {
		t.neutral.AlertTmpl.MessageID = msgID
	}
	if strings.Contains(content, t.enemies.AlertTmpl.Header) {
		t.enemies.AlertTmpl.MessageID = msgID
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
