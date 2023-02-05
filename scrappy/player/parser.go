package player

import (
	"darkbot/scrappy/shared/records"
	"darkbot/utils/logger"
	"encoding/json"
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

func (b playerParser) Parse(body []byte) records.StampedObjects[Player] {
	record := records.StampedObjects[Player]{}.New()

	playerData := SerializedPlayers{}
	if err := json.Unmarshal(body, &playerData); err != nil {
		logger.CheckPanic(err, "unable to unmarshal base request")
	}

	for _, serializedPlayer := range playerData.Players {
		record.Add(
			serializedPlayer.Name,
			Player{
				Time:   serializedPlayer.Time,
				Name:   serializedPlayer.Name,
				System: serializedPlayer.System,
				Region: serializedPlayer.Region,
			},
		)
	}
	return record
}
