package eventview

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/viewer/apis"
	"darkbot/app/viewer/views"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type EventRenderer struct {
	main views.TemplateShared
	api  *apis.API
}

func NewEventRenderer(api *apis.API) *EventRenderer {
	base := EventRenderer{}
	base.api = api
	base.main.Header = "#darkbot-event-view"
	return &base
}

func (t *EventRenderer) DiscoverMessageID(content string, msgID types.DiscordMessageID) {
	if strings.Contains(content, t.main.Header) {
		t.main.MessageID = msgID
	}
}

func (t *EventRenderer) MatchMessageID(messageID types.DiscordMessageID) bool {
	return messageID == t.main.MessageID
}

func (t *EventRenderer) Send() {
	t.main.Send(t.api)
}

type PlayerTemplate struct {
	Time   string
	Name   string
	System string
}

func (t *EventRenderer) Render() error {
	record, err := t.api.Scrappy.GetPlayerStorage().GetLatestRecord()
	if logus.CheckWarn(err, "unable to get players") {
		return err
	}

	logus.Debug("rendered events", logus.DiscordMessageID(t.main.MessageID))

	eventTags, err := t.api.Players.Events.TagsList(t.api.ChannelID)
	logus.CheckWarn(err, "failed to acquire player event list", logus.ChannelID(t.api.ChannelID))

	if len(eventTags) > 0 {
		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("**%s** %s\n", t.main.Header, time.Now().String()))
		sb.WriteString("**Event table of players**\n")
		sb.WriteString("```json\n")

		for _, eventTag := range eventTags {
			sb.WriteString(fmt.Sprintf(`"%s": `, eventTag))

			matchedPlayers := []PlayerTemplate{}
			for _, player := range record.List {
				if views.TagContains(player.Name, []types.Tag{eventTag}) {
					matchedPlayers = append(matchedPlayers, PlayerTemplate{Name: player.Name, Time: player.Time, System: player.System})
					continue
				}
			}
			result, err := json.Marshal(matchedPlayers)
			logus.CheckError(err, "failed to marshal event matched players")
			sb.WriteString(fmt.Sprintf("%s", string(result)))

			sb.WriteString("\n")
		}

		sb.WriteString("```\n")
		t.main.Content = sb.String()

	}

	return nil
}
