package player

import (
	"darkbot/scrappy/shared/records"
	"darkbot/settings/utils/logger"
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

func (b playerParser) Parse(body []byte) (records.StampedObjects[Player], error) {
	record := records.StampedObjects[Player]{}.New()

	playerData := SerializedPlayers{}
	if err := json.Unmarshal(body, &playerData); err != nil {
		logger.CheckWarn(err, "unable to unmarshal player request")
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
