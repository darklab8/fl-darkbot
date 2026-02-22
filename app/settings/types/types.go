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

func GetI(pointer *int, defaul int) int {
	if pointer == nil {
		return 0
	}
	return *pointer
}
func GetF(pointer *float64, defaul float64) float64 {
	if pointer == nil {
		return 0
	}
	return *pointer
}
func GetS(pointer *string, defaul string) string {
	if pointer == nil {
		return ""
	}
	return *pointer
}
