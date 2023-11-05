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
)

type EventView struct {
	main views.ViewTable

	*views.SharedViewTableSplitter
}

func NewEventRenderer(api *apis.API) *EventView {
	base := EventView{}
	base.main.ViewID = "#darkbot-event-view"

	base.SharedViewTableSplitter = views.NewSharedViewSplitter(api, &base, &base.main)

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

	// logus.Debug("rendered events", logus.DiscordMessageID(t.main.MessageID)) // TODO delete as u finished refactor

	eventTags, err := t.GetAPI().Players.Events.TagsList(t.GetAPI().ChannelID)
	logus.CheckWarn(err, "failed to acquire player event list", logus.ChannelID(t.GetAPI().ChannelID))

	if len(eventTags) > 0 {
		var beginning strings.Builder
		var end strings.Builder

		// Looks like identical :thinking: // TODO delete as u finished refactor
		// sb.WriteString(fmt.Sprintf("**%s** %s\n", t.main.ViewID, time.Now().String()))

		beginning.WriteString("**Event table of players**\n")
		beginning.WriteString("```json\n")
		t.main.ViewBeginning = types.ViewBeginning(beginning.String())

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
			record.WriteString(fmt.Sprintf("%s", string(result)))

			record.WriteString("\n")
			t.main.AppendRecord(types.ViewRecord(record.String()))
		}

		end.WriteString("```\n")
		t.main.ViewEnd = types.ViewEnd(end.String())

	}

	return nil
}
