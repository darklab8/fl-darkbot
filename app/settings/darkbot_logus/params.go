package darkbot_logus

import (
	"darkbot/app/forumer/forum_types"
	"darkbot/app/settings/types"
	"fmt"

	"github.com/darklab8/darklab_goutils/goutils/logus"
	"gorm.io/gorm"
)

func Records[T any](value []T) logus.SlogParam {
	return logus.Items[T](value, "records")
}

func Args(value []string) logus.SlogParam {
	return logus.Items[string](value, "args")
}

func Tags(value []types.Tag) logus.SlogParam {
	return logus.Items[types.Tag](value, "tags")
}

func APIUrl(value types.APIurl) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["api_url"] = string(value)
		c.Params["records_len"] = fmt.Sprintf("%d", len(value))
	}
}

func ScrappyLoopDelay(value types.ScrappyLoopDelay) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["loop_delay"] = fmt.Sprintf("%d", value)
	}
}

func ChannelID(value types.DiscordChannelID) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["channel_id"] = string(value)
	}
}

func ChannelIDs(value []types.DiscordChannelID) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["channel_ids"] = fmt.Sprintf("%v", value)
	}
}

func MessageID(value types.DiscordMessageID) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["message_id"] = string(value)
	}
}

func OwnerID(value types.DiscordOwnerID) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["owner_id"] = string(value)
	}
}

func Body(value []byte) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["body"] = string(value)
	}
}

func PingMessage(value types.PingMessage) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["ping_message"] = string(value)
	}
}

func ErrorMsg(value string) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["error_message"] = string(value)
	}
}

func Dbpath(value types.Dbpath) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["db_path"] = string(value)
	}
}

func Tag(value types.Tag) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["tag"] = string(value)
	}
}

func GormResult(result *gorm.DB) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["result.rows_affected"] = fmt.Sprintf("%d", result.RowsAffected)
		c.Params["result.error_msg"] = fmt.Sprintf("%v", result.Error)
		c.Params["result.error_type"] = fmt.Sprintf("%T", result.Error)
	}
}

func DiscordMessageID(value types.DiscordMessageID) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["discord_msg_id"] = string(value)
	}
}

func Thread(value *forum_types.LatestThread) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["thread_name"] = string(value.ThreadShortName)
		c.Params["thread_link"] = string(value.ThreadLink)
		c.Params["thread_id"] = string(value.ThreadID)
	}
}

func DiscordMessage(value string) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["discord_message"] = value
	}
}

func Post(value *forum_types.Post) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["post_id"] = string(value.PostID)
		c.Params["post_author_name"] = string(value.PostAuthorName)
		c.Params["post_thread_id"] = string(value.ThreadID)
		c.Params["post_thread_full_name"] = string(value.ThreadFullName)
	}
}
