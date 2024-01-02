package views

import (
	"darkbot/app/settings/darkbot_logus"
	"darkbot/app/settings/types"
	"darkbot/app/viewer/apis"
	"darkbot/app/viewer/views/viewer_msg"
	"fmt"
	"strings"

	"github.com/darklab8/darklab_goutils/goutils/utils"
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

func (t *SharedViewTableSplitter) RenderView() error {
	err := t.original_view.GenerateRecords()

	for _, view := range t.views {
		msg_count := 0
		msg := viewer_msg.NewMsg(view.msgShared, msg_count, t.channelID)

		for _, record := range view.viewRecords {
			if len(*record)+msg.Len() > DiscordMsgLimit {
				view.msgs = append(view.msgs, msg)
				msg_count += 1
				msg = viewer_msg.NewMsg(view.msgShared, msg_count, t.channelID)
			}

			msg.AppendRecordToMsg(record)
		}

		if msg.HasRecords() {
			view.msgs = append(view.msgs, msg)
		}
	}

	return err
}

// Time comlexity: Must be called only after Generate()
func (t *SharedViewTableSplitter) DiscoverMessageID(content string, msgID types.DiscordMessageID) {
	for _, view := range t.views {
		for _, msg := range view.msgs {
			if strings.Contains(content, string(msg.GetViewEnumeratedID())) {
				darkbot_logus.Log.Debug(fmt.Sprintf("discovered content to ViewEnumeratedID=%v", viewer_msg.LogusMsg(msg)))
				msg.SetMessageID(msgID)
			}
		}
	}
}

// Time Complexity: Must be called only after DiscoverMessageID
func (t *SharedViewTableSplitter) MatchMessageID(messageID types.DiscordMessageID) bool {
	for _, view := range t.views {
		for _, msg := range view.msgs {
			if msg.GetMessageID() == messageID {
				darkbot_logus.Log.Debug(fmt.Sprintf("found match messageID=%v to msg.MessageID=%v", messageID, msg.GetMessageID()))
				return true
			}
		}
	}
	return false
}

func (t *SharedViewTableSplitter) Send() {
	utils.TimeMeasure(func() {
		for _, view := range t.views {
			utils.TimeMeasure(func() {
				for _, msg := range view.msgs {
					utils.TimeMeasure(func() {
						msg.Send(t.api)
					}, "SharedViewTableSplitter.send", viewer_msg.LogusMsg(msg))
				}
			}, fmt.Sprintf("SharedViewTableSplitter.send.view=%v", view))
		}
	}, fmt.Sprintf("SharedViewTableSplitter.send=%v", t))
}
