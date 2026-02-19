package types

type Dbpath string

type APIurl string

type ScrappyLoopDelay int
type ViewerLoopDelay int

type ViewerDelayBetweenChannels int

type DiscordChannelID string

type DiscordMessageID string

type DiscordOwnerID string

type PingMessage string

type Tag string

type ViewID string
type ViewEnumeratedID string

// Message
type ViewTimeStamp string
type ViewBeginning string
type ViewHeader string
type ViewRecord string
type ViewEnd string

type OrderKey string

func GetF(pointer *float64) float64 {
	if pointer == nil {
		return 0
	}
	return *pointer
}
func GetS(pointer *string) string {
	if pointer == nil {
		return ""
	}
	return *pointer
}
