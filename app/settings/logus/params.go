package logus

import (
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

type slogGroup struct {
	params map[string]string
}

func (s slogGroup) Render() slog.Attr {
	anies := []any{}
	for key, value := range s.params {
		anies = append(anies, key)
		anies = append(anies, value)
	}

	return slog.Group("extras", anies...)
}

type slogParam func(r *slogGroup)

func newSlogGroup(opts ...slogParam) slog.Attr {
	client := &slogGroup{params: make(map[string]string)}
	for _, opt := range opts {
		opt(client)
	}

	return (*client).Render()
}

func TestParam(value int) slogParam {
	return func(c *slogGroup) {
		c.params["test_param"] = fmt.Sprintf("%d", value)
	}
}

func Expected(value any) slogParam {
	return func(c *slogGroup) {
		c.params["expected"] = fmt.Sprintf("%v", value)
	}
}
func Actual(value any) slogParam {
	return func(c *slogGroup) {
		c.params["actual"] = fmt.Sprintf("%v", value)
	}
}

func OptError(err error) slogParam {
	return func(c *slogGroup) {
		c.params["error_msg"] = fmt.Sprintf("%v", err)
		c.params["error_type"] = fmt.Sprintf("%T", err)
	}
}

func FilePath(value string) slogParam {
	return func(c *slogGroup) {
		c.params["filepath"] = fmt.Sprintf("%v", value)
	}
}

func Regex(value types.RegExp) slogParam {
	return func(c *slogGroup) {
		c.params["regexp"] = fmt.Sprintf("%v", value)
	}
}

func Items[T any](value []T, item_name string) slogParam {
	return func(c *slogGroup) {
		c.params[item_name] = fmt.Sprintf("%v", value)
		c.params[fmt.Sprintf("%s_len", item_name)] = fmt.Sprintf("%d", len(value))
	}
}

func Records[T any](value []T) slogParam {
	return Items[T](value, "records")
}

func Args(value []string) slogParam {
	return Items[string](value, "args")
}

func Tags(value []types.Tag) slogParam {
	return Items[types.Tag](value, "tags")
}

func APIUrl(value types.APIurl) slogParam {
	return func(c *slogGroup) {
		c.params["api_url"] = string(value)
		c.params["records_len"] = fmt.Sprintf("%d", len(value))
	}
}

func ScrappyLoopDelay(value types.ScrappyLoopDelay) slogParam {
	return func(c *slogGroup) {
		c.params["loop_delay"] = fmt.Sprintf("%d", value)
	}
}

func ChannelID(value types.DiscordChannelID) slogParam {
	return func(c *slogGroup) {
		c.params["channel_id"] = string(value)
	}
}

func ChannelIDs(value []types.DiscordChannelID) slogParam {
	return func(c *slogGroup) {
		c.params["channel_ids"] = fmt.Sprintf("%v", value)
	}
}

func MessageID(value types.DiscordMessageID) slogParam {
	return func(c *slogGroup) {
		c.params["message_id"] = string(value)
	}
}

func OwnerID(value types.DiscordOwnerID) slogParam {
	return func(c *slogGroup) {
		c.params["owner_id"] = string(value)
	}
}

func Body(value []byte) slogParam {
	return func(c *slogGroup) {
		c.params["body"] = string(value)
	}
}

func PingMessage(value types.PingMessage) slogParam {
	return func(c *slogGroup) {
		c.params["ping_message"] = string(value)
	}
}

func ErrorMsg(value string) slogParam {
	return func(c *slogGroup) {
		c.params["error_message"] = string(value)
	}
}

func Dbpath(value types.Dbpath) slogParam {
	return func(c *slogGroup) {
		c.params["db_path"] = string(value)
	}
}

func Tag(value types.Tag) slogParam {
	return func(c *slogGroup) {
		c.params["tag"] = string(value)
	}
}

func GormResult(result *gorm.DB) slogParam {
	return func(c *slogGroup) {
		c.params["result.rows_affected"] = fmt.Sprintf("%d", result.RowsAffected)
		c.params["result.error_msg"] = fmt.Sprintf("%v", result.Error)
		c.params["result.error_type"] = fmt.Sprintf("%T", result.Error)
	}
}
