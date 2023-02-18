package templ

import (
	"darkbot/dtypes"
	"darkbot/scrappy/player"
	"darkbot/utils"
	"darkbot/utils/logger"
	"darkbot/viewer/apis"
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
	TemplateShared
}
type PlayersEnemies struct {
	TemplateShared
}
type PlayersNeutral struct {
	TemplateShared
}

type PlayersTemplates struct {
	friends PlayersFriends
	neutral PlayersNeutral
	enemies PlayersEnemies
}

func NewTemplatePlayers(channelID string, dbpath dtypes.Dbpath) PlayersTemplates {
	templator := PlayersTemplates{}
	templator.friends.API = apis.NewAPI(channelID, dbpath)
	templator.neutral.API = apis.NewAPI(channelID, dbpath)
	templator.enemies.API = apis.NewAPI(channelID, dbpath)
	templator.friends.Header = "#darkbot-players-friends"
	templator.neutral.Header = "#darkbot-players-neutral"
	templator.enemies.Header = "#darkbot-players-enemies"
	return templator
}

func (b *PlayersTemplates) Setup(channelID string) {
	b.neutral.MessageID = ""
	b.enemies.MessageID = ""
	b.friends.MessageID = ""
	b.neutral.Content = ""
	b.enemies.Content = ""
	b.friends.Content = ""
	b.neutral.API.ChannelID = channelID
	b.enemies.API.ChannelID = channelID
	b.friends.API.ChannelID = channelID
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
	record, err := t.friends.Scrappy.PlayerStorage.GetLatestRecord()
	if err != nil {
		return
	}

	systemTags, _ := t.neutral.Systems.TagsList(t.neutral.ChannelID)
	regionTags, _ := t.neutral.Regions.TagsList(t.neutral.ChannelID)
	friendTags, _ := t.neutral.Friends.TagsList(t.neutral.ChannelID)
	enemyTags, _ := t.neutral.Enemies.TagsList(t.neutral.ChannelID)
	fmt.Println("systemTags=", systemTags)
	fmt.Println("regionTags=", regionTags)
	fmt.Println("friendTags=", friendTags)
	fmt.Println("enemyTags=", enemyTags)
	fmt.Println("record.List=", record.List)

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

	logger.Debug("friendPlayers=", friendPlayers)
	logger.Debug("enemyPlayers=", enemyPlayers)
	logger.Debug("neutralPlayers=", neutralPlayers)

	if len(neutralPlayers) > 0 {
		t.neutral.Content = utils.TmpRender(playerTemplate, TemplateRendrerPlayerInput{
			Header:      t.neutral.Header,
			LastUpdated: time.Now().String(),
			Players:     neutralPlayers,
			TableName:   "**Neutral players in tracked systems and regions**",
		})
	}

	if len(enemyPlayers) > 0 {
		t.enemies.Content = utils.TmpRender(playerTemplate, TemplateRendrerPlayerInput{
			Header:      t.enemies.Header,
			LastUpdated: time.Now().String(),
			Players:     enemyPlayers,
			TableName:   "**Enemy players in tracked systems and regions**",
		})
	}

	if len(friendPlayers) > 0 {
		t.friends.Content = utils.TmpRender(playerTemplate, TemplateRendrerPlayerInput{
			Header:      t.friends.Header,
			LastUpdated: time.Now().String(),
			Players:     friendPlayers,
			TableName:   "**Friend players in all systems and regions**",
		})
	}
}

func (t *PlayersTemplates) Send() {
	t.friends.Send()
	t.neutral.Send()
	t.enemies.Send()
}

func (t *PlayersTemplates) MatchMessageID(messageID string) bool {

	if messageID == t.friends.MessageID {
		return true
	}
	if messageID == t.neutral.MessageID {
		return true
	}
	if messageID == t.enemies.MessageID {
		return true
	}
	return false
}

func (t *PlayersTemplates) DiscoverMessageID(content string, msgID string) {
	if strings.Contains(content, t.friends.Header) {
		t.friends.MessageID = msgID
	}
	if strings.Contains(content, t.neutral.Header) {
		t.neutral.MessageID = msgID
	}
	if strings.Contains(content, t.enemies.Header) {
		t.enemies.MessageID = msgID
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
