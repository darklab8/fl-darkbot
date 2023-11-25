package eventview

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/viewer/apis"
	"darkbot/app/viewer/views"
	"darkbot/app/viewer/views/viewer_msg"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

type EventView struct {
	main *views.ViewTable

	*views.SharedViewTableSplitter
	channelID types.DiscordChannelID
}

func NewEventRenderer(api *apis.API, channelID types.DiscordChannelID) *EventView {
	base := EventView{}
	base.main = views.NewViewTable(viewer_msg.NewTableMsg(
		types.ViewID("#darkbot-event-view"),
		types.ViewHeader("**Event table of players**\n"),
		types.ViewBeginning("```json\n"),
		types.ViewEnd("```\n"),
	))
	base.channelID = channelID

	base.SharedViewTableSplitter = views.NewSharedViewSplitter(api, channelID, &base, base.main)

	return &base
}

type PlayerTemplate struct {
	Time   string
	Name   string
	System string
}

func (t *EventView) GenerateRecords() error {
	player_record, err := t.GetAPI().Scrappy.GetPlayerStorage().GetLatestRecord()
	if logus.CheckWarn(err, "unable to get players") {
		return err
	}

	eventTags, err := t.GetAPI().Players.Events.TagsList(t.channelID)
	logus.CheckDebug(err, "failed to acquire player event list", logus.ChannelID(t.channelID))

	if len(eventTags) > 0 {
		for _, eventTag := range eventTags {
			var record strings.Builder
			record.WriteString(fmt.Sprintf(`"%s": `, eventTag))

			matchedPlayers := []PlayerTemplate{}
			for _, player := range player_record.List {
				if views.TagContains(player.Name, []types.Tag{eventTag}) {
					matchedPlayers = append(matchedPlayers, PlayerTemplate{Name: player.Name, Time: player.Time, System: player.System})
					continue
				}
			}
			result, err := json.Marshal(matchedPlayers)
			logus.CheckError(err, "failed to marshal event matched players")
			record.WriteString(string(result))

			record.WriteString("\n")
			t.main.AppendRecord(types.ViewRecord(record.String()))
		}
	}

	return nil
}
