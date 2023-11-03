package consoler_types

import "darkbot/app/settings/types"

type ChannelParams struct {
	channelID types.DiscordChannelID
	dbpath    types.Dbpath
}

func (c *ChannelParams) SetChannelID(channelID types.DiscordChannelID) {
	c.channelID = channelID
}

func NewChannelParams(
	channelID types.DiscordChannelID,
	dbpath types.Dbpath,

) *ChannelParams {
	return &ChannelParams{channelID: channelID, dbpath: dbpath}
}

func (c ChannelParams) GetChannelID() types.DiscordChannelID { return c.channelID }
func (c ChannelParams) GetDbpath() types.Dbpath              { return c.dbpath }
