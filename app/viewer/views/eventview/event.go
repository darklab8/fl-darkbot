package eventview

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/darklab/fl-darkbot/app/settings/logus"
	"github.com/darklab/fl-darkbot/app/settings/types"
	"github.com/darklab/fl-darkbot/app/viewer/apis"
	"github.com/darklab/fl-darkbot/app/viewer/views"
	"github.com/darklab/fl-darkbot/app/viewer/views/viewer_msg"
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
	if logus.Log.CheckWarn(err, "unable to get players") {
		return err
	}

	eventTags, err := t.GetAPI().Players.Events.TagsList(t.channelID)
	logus.Log.CheckDebug(err, "failed to acquire player event list", logus.ChannelID(t.channelID))

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
			logus.Log.CheckError(err, "failed to marshal event matched players")
			record.WriteString(string(result))

			record.WriteString("\n")
			t.main.AppendRecord(types.ViewRecord(record.String()))
		}
	}

	return nil
}
