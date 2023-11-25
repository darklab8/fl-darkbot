package views

import (
	"darkbot/app/settings/types"
	"darkbot/app/viewer/views/viewer_msg"
)

type ViewTable struct {
	msgShared   *viewer_msg.MsgShared
	msgs        []*viewer_msg.Msg
	viewRecords []*types.ViewRecord
}

func NewViewTable(msgShared *viewer_msg.MsgShared) *ViewTable {
	return &ViewTable{
		msgShared: msgShared,
	}
}

func (v *ViewTable) HasRecords() bool { return len(v.viewRecords) > 0 }

func (m *ViewTable) RecordCount() int { return len(m.viewRecords) }

func (v *ViewTable) GetMsgs() []*viewer_msg.Msg { return v.msgs }

func (v *ViewTable) SetHeader(header types.ViewHeader) { v.msgShared.SetHeader(header) }

func (v *ViewTable) AppendRecord(record types.ViewRecord) {
	v.viewRecords = append(v.viewRecords, &record)
}

type ViewHeader string

const DiscordMsgLimit = 2000 - 150

type View interface {
	DiscoverMessageID(content string, msgID types.DiscordMessageID)
	RenderView() error
	Send()
	MatchMessageID(messageID types.DiscordMessageID) bool
}
