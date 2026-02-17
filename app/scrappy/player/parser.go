package player

import (
	"encoding/json"

	"github.com/darklab8/fl-darkbot/app/scrappy/shared/records"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
)

type SerializedPlayers struct {
	Error   interface{} `json:"error"`
	Players []struct {
		Time   string `json:"time"`
		Name   string `json:"name"`
		System string `json:"system"`
		Region string `json:"region"`
	} `json:"players"`
	Timestamp string `json:"timestamp"`
}

type playerParser struct {
}

func (b playerParser) Parse(body []byte) (records.StampedObjects[Player], error) {
	record := records.NewStampedObjects[Player]()

	playerData := SerializedPlayers{}
	if err := json.Unmarshal(body, &playerData); err != nil {
		logus.Log.CheckWarn(err, "unable to unmarshal player request")
		return record, err
	}

	for _, serializedPlayer := range playerData.Players {
		record.Add(
			Player{
				Time:   serializedPlayer.Time,
				Name:   serializedPlayer.Name,
				System: serializedPlayer.System,
				Region: serializedPlayer.Region,
			},
		)
	}
	return record, nil
}
