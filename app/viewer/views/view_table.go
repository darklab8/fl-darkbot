package views

import (
	"darkbot/app/settings/types"
	"fmt"
	"time"
)

type MsgAction int

const (
	ActSend MsgAction = iota
	ActEdit
)

type ViewTable struct {
	MsgShared
	Msgs        []*Msg
	ViewRecords []*types.ViewRecord
}

func (v *ViewTable) AppendRecord(record types.ViewRecord) {
	v.ViewRecords = append(v.ViewRecords, &record)
}

type ViewHeader string

func (t ViewTable) GeneratedHeader() ViewHeader {
	return ViewHeader(fmt.Sprintf("**%s** %s\n", t.ViewID, time.Now().String()))
}

const DiscordMsgLimit = 2000 - 50

func (t ViewTable) GeneratedViews(returned_view func()) {
	returned_view()
}

type View interface {
	DiscoverMessageID(content string, msgID types.DiscordMessageID)
	Render() error
	Send()
	MatchMessageID(messageID types.DiscordMessageID) bool
}
