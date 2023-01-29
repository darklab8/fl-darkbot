package templ

import (
	"darkbot/viewer/apis"
	"strings"
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
type PlayersAll struct {
	TemplateShared
}

type PlayersTemplates struct {
	friends PlayersFriends
	all     PlayersAll
	enemies PlayersEnemies
}

func NewTemplatePlayers(channelID string) *PlayersTemplates {
	templator := PlayersTemplates{}
	templator.friends.API = apis.NewAPI(channelID)
	templator.all.API = apis.NewAPI(channelID)
	templator.enemies.API = apis.NewAPI(channelID)
	templator.friends.Header = "#darkbot-players-friends"
	templator.all.Header = "#darkbot-players-all"
	templator.enemies.Header = "#darkbot-players-enemies"
	return &templator
}

func (t *PlayersTemplates) Render() {

}

func (t *PlayersTemplates) Send() {
	t.friends.Send()
	t.all.Send()
	t.enemies.Send()
}

func (t *PlayersTemplates) MatchMessageID(messageID string) bool {

	if messageID == t.friends.MessageID {
		return true
	}
	if messageID == t.all.MessageID {
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
	if strings.Contains(content, t.all.Header) {
		t.all.MessageID = msgID
	}
	if strings.Contains(content, t.enemies.Header) {
		t.enemies.MessageID = msgID
	}
}
