package printer

import "darkbot/settings/types"

type ChannelInfo struct {
	ChannelID types.DiscordChannelID
	Dbpath    types.Dbpath
}

func NewChannelInfo(
	ChannelID types.DiscordChannelID,
	Dbpath types.Dbpath,
) ChannelInfo {
	return ChannelInfo{ChannelID: ChannelID, Dbpath: Dbpath}
}
