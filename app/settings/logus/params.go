package logus

import (
	"darkbot/app/forumer/forum_types"
	"darkbot/app/settings/types"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

func logGroupFiles() slog.Attr {
	return slog.Group("files",
		"file3", GetCallingFile(3),
		"file4", GetCallingFile(4),
	)
}

type SlogGroup struct {
	Params map[string]string
}

func (s SlogGroup) Render() slog.Attr {
	anies := []any{}
	for key, value := range s.Params {
		anies = append(anies, key)
		anies = append(anies, value)
	}

	return slog.Group("extras", anies...)
}

type SlogParam func(r *SlogGroup)

func newSlogGroup(opts ...SlogParam) slog.Attr {
	client := &SlogGroup{Params: make(map[string]string)}
	for _, opt := range opts {
		opt(client)
	}

	return (*client).Render()
}

func TestParam(value int) SlogParam {
	return func(c *SlogGroup) {
		c.Params["test_param"] = fmt.Sprintf("%d", value)
	}
}

func Expected(value any) SlogParam {
	return func(c *SlogGroup) {
		c.Params["expected"] = fmt.Sprintf("%v", value)
	}
}
func Actual(value any) SlogParam {
	return func(c *SlogGroup) {
		c.Params["actual"] = fmt.Sprintf("%v", value)
	}
}

func OptError(err error) SlogParam {
	return func(c *SlogGroup) {
		c.Params["error_msg"] = fmt.Sprintf("%v", err)
		c.Params["error_type"] = fmt.Sprintf("%T", err)
	}
}

func FilePath(value string) SlogParam {
	return func(c *SlogGroup) {
		c.Params["filepath"] = fmt.Sprintf("%v", value)
	}
}

func Regex(value types.RegExp) SlogParam {
	return func(c *SlogGroup) {
		c.Params["regexp"] = fmt.Sprintf("%v", value)
	}
}

func Items[T any](value []T, item_name string) SlogParam {
	return func(c *SlogGroup) {
		sliced_string := fmt.Sprintf("%v", value)
		if len(sliced_string) > 300 {
			sliced_string = sliced_string[:300] + "...sliced string"
		}
		c.Params[item_name] = sliced_string
		c.Params[fmt.Sprintf("%s_len", item_name)] = fmt.Sprintf("%d", len(value))
	}
}

func Records[T any](value []T) SlogParam {
	return Items[T](value, "records")
}

func Args(value []string) SlogParam {
	return Items[string](value, "args")
}

func Tags(value []types.Tag) SlogParam {
	return Items[types.Tag](value, "tags")
}

func APIUrl(value types.APIurl) SlogParam {
	return func(c *SlogGroup) {
		c.Params["api_url"] = string(value)
		c.Params["records_len"] = fmt.Sprintf("%d", len(value))
	}
}

func ScrappyLoopDelay(value types.ScrappyLoopDelay) SlogParam {
	return func(c *SlogGroup) {
		c.Params["loop_delay"] = fmt.Sprintf("%d", value)
	}
}

func ChannelID(value types.DiscordChannelID) SlogParam {
	return func(c *SlogGroup) {
		c.Params["channel_id"] = string(value)
	}
}

func ChannelIDs(value []types.DiscordChannelID) SlogParam {
	return func(c *SlogGroup) {
		c.Params["channel_ids"] = fmt.Sprintf("%v", value)
	}
}

func MessageID(value types.DiscordMessageID) SlogParam {
	return func(c *SlogGroup) {
		c.Params["message_id"] = string(value)
	}
}

func OwnerID(value types.DiscordOwnerID) SlogParam {
	return func(c *SlogGroup) {
		c.Params["owner_id"] = string(value)
	}
}

func Body(value []byte) SlogParam {
	return func(c *SlogGroup) {
		c.Params["body"] = string(value)
	}
}

func PingMessage(value types.PingMessage) SlogParam {
	return func(c *SlogGroup) {
		c.Params["ping_message"] = string(value)
	}
}

func ErrorMsg(value string) SlogParam {
	return func(c *SlogGroup) {
		c.Params["error_message"] = string(value)
	}
}

func Dbpath(value types.Dbpath) SlogParam {
	return func(c *SlogGroup) {
		c.Params["db_path"] = string(value)
	}
}

func Tag(value types.Tag) SlogParam {
	return func(c *SlogGroup) {
		c.Params["tag"] = string(value)
	}
}

func GormResult(result *gorm.DB) SlogParam {
	return func(c *SlogGroup) {
		c.Params["result.rows_affected"] = fmt.Sprintf("%d", result.RowsAffected)
		c.Params["result.error_msg"] = fmt.Sprintf("%v", result.Error)
		c.Params["result.error_type"] = fmt.Sprintf("%T", result.Error)
	}
}

func DiscordMessageID(value types.DiscordMessageID) SlogParam {
	return func(c *SlogGroup) {
		c.Params["discord_msg_id"] = string(value)
	}
}

func Thread(value *forum_types.LatestThread) SlogParam {
	return func(c *SlogGroup) {
		c.Params["thread_name"] = string(value.ThreadShortName)
		c.Params["thread_link"] = string(value.ThreadLink)
		c.Params["thread_id"] = string(value.ThreadID)
	}
}

func DiscordMessage(value string) SlogParam {
	return func(c *SlogGroup) {
		c.Params["discord_message"] = value
	}
}

func Post(value *forum_types.Post) SlogParam {
	return func(c *SlogGroup) {
		c.Params["post_id"] = string(value.PostID)
		c.Params["post_author_name"] = string(value.PostAuthorName)
		c.Params["post_thread_id"] = string(value.ThreadID)
		c.Params["post_thread_full_name"] = string(value.ThreadFullName)
	}
}
