package logus

import (
	"fmt"
	"log/slog"

	"github.com/darklab8/fl-darkbot/app/forumer/forum_types"
	"github.com/darklab8/fl-darkbot/app/settings/types"

	"github.com/darklab8/go-typelog/typelog"
	"gorm.io/gorm"
)

func Records[T any](value []T) typelog.LogType {
	return typelog.Items[T]("records", value)
}

func Args(value []string) typelog.LogType {
	return typelog.Items[string]("args", value)
}

func Tags(value []types.Tag) typelog.LogType {
	return typelog.Items[types.Tag]("tags", value)
}

func APIUrl(value types.APIurl) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(
			slog.String("api_url", string(value)),
			slog.Int("records_len", len(value)),
		)
	}
}

func ScrappyLoopDelay(value types.ScrappyLoopDelay) typelog.LogType {
	return typelog.Int("loop_delay", int(value))
}

func ChannelID(value types.DiscordChannelID) typelog.LogType {
	return typelog.String("channel_id", string(value))
}

func MsgContent(value string) typelog.LogType {
	return typelog.String("content", ShortenedMsg(value))
}

func ShortenedMsg(msg string) string {
	if len(msg) > 100 {
		return msg[:100]
	}
	return msg
}

func ChannelIDs(value []types.DiscordChannelID) typelog.LogType {
	return typelog.Any("channel_ids", value)
}

func MessageID(value types.DiscordMessageID) typelog.LogType {
	return typelog.String("message_id", string(value))
}

func OwnerID(value types.DiscordOwnerID) typelog.LogType {
	return typelog.String("owner_id", string(value))
}

func Body(value []byte) typelog.LogType {
	return typelog.Bytes("body", value)
}

func PingMessage(value types.PingMessage) typelog.LogType {
	return typelog.String("ping_message", string(value))
}

func ErrorMsg(value string) typelog.LogType {
	return typelog.String("error_message", value)
}

func Dbpath(value types.Dbpath) typelog.LogType {
	return typelog.String("db_path", string(value))
}

func Tag(value types.Tag) typelog.LogType {
	return typelog.String("tag", string(value))
}

func GormResult(result *gorm.DB) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(
			slog.Int64("result.rows_affected", result.RowsAffected),
			slog.Any("result.error_msg", result.Error),
			slog.String("result.error_type", fmt.Sprintf("%T", result.Error)),
		)
	}
}

func DiscordMessageID(value types.DiscordMessageID) typelog.LogType {
	return typelog.String("discord_msg_id", string(value))
}

func Thread(value *forum_types.LatestThread) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(
			slog.String("thread_name", string(value.ThreadShortName)),
			slog.String("thread_link", string(value.ThreadLink)),
			slog.String("thread_id", string(value.ThreadID)),
		)
	}
}

func DiscordMessage(value string) typelog.LogType {
	return typelog.String("discord_message", value)
}

func Post(value *forum_types.Post) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(
			slog.String("post_id", string(value.PostID)),
			slog.String("post_author_name", string(value.PostAuthorName)),
			slog.String("post_thread_id", string(value.ThreadID)),
			slog.String("post_thread_full_name", string(value.ThreadFullName)),
		)
	}
}
