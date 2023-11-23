package views

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/viewer/apis"
	"fmt"
	"strings"
)

type SharedViewTableSplitter struct {
	views         []*ViewTable
	api           *apis.API
	channelID     types.DiscordChannelID
	original_view OriginalRenderer
}

func NewSharedViewSplitter(
	api *apis.API,
	channelID types.DiscordChannelID,
	original_view OriginalRenderer,
	sh_templates ...*ViewTable,
) *SharedViewTableSplitter {
	s := &SharedViewTableSplitter{
		api:           api,
		original_view: original_view,
		channelID:     channelID,
	}

	s.views = append(s.views, sh_templates...)
	return s
}

func (s *SharedViewTableSplitter) GetAPI() *apis.API { return s.api }

type OriginalRenderer interface {
	GenerateRecords() error
}

func (t *SharedViewTableSplitter) Render() error {
	err := t.original_view.GenerateRecords()

	for _, view := range t.views {
		msg_count := 0
		msg := NewMsg(view.MsgShared, msg_count, t.channelID)

		for _, record := range view.ViewRecords {
			if len(*record)+msg.Len() > DiscordMsgLimit {
				view.Msgs = append(view.Msgs, msg)
				msg_count += 1
				msg = NewMsg(view.MsgShared, msg_count, t.channelID)
			}

			msg.Records = append(msg.Records, record)
		}

		if len(msg.Records) > 0 {
			view.Msgs = append(view.Msgs, msg)
		}
	}

	return err
}

// Time comlexity: Must be called only after Generate()
func (t *SharedViewTableSplitter) DiscoverMessageID(content string, msgID types.DiscordMessageID) {
	for _, view := range t.views {
		for _, msg := range view.Msgs {
			if strings.Contains(content, string(msg.ViewEnumeratedID)) {
				logus.Debug(fmt.Sprintf("discovered content to ViewEnumeratedID=%v", msg.ViewEnumeratedID))
				msg.MessageID = msgID
			}
		}
	}
}

// Time Complexity: Must be called only after DiscoverMessageID
func (t *SharedViewTableSplitter) MatchMessageID(messageID types.DiscordMessageID) bool {
	for _, view := range t.views {
		for _, msg := range view.Msgs {
			if msg.MessageID == messageID {
				logus.Debug(fmt.Sprintf("found match messageID=%v to msg.MessageID=%v", messageID, msg.MessageID))
				return true
			}
		}
	}
	return false
}

func (t *SharedViewTableSplitter) Send() {
	for _, view := range t.views {
		for _, msg := range view.Msgs {
			msg.Send(t.api)
		}
	}
}
