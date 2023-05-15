package redis

import (
	"context"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/samber/lo"
)

type Channel struct {
	re *Redis
}

func (c *Channel) Reset() error {
	return c.re.PatternDel(makeKey("channels", discord.NullChannelID))
}

func (c *Channel) key(channelID discord.ChannelID) string {
	return makeKey("channels", channelID)
}

func (c *Channel) Channel(id discord.ChannelID) (*discord.Channel, error) {
	return Get[*discord.Channel](c.re, c.key(id))
}

func (c *Channel) CreatePrivateChannel(recipient discord.UserID) (*discord.Channel, error) {
	channelID, err := Get[discord.ChannelID](c.re, "channels:private:"+recipient.String())

	if err != nil {
		return nil, err
	}

	return Get[*discord.Channel](c.re, c.key(channelID))
}

func (c *Channel) Channels(guildID discord.GuildID) ([]discord.Channel, error) {
	cmd := c.re.client.SMembers(context.Background(), "guilds:"+guildID.String()+":channels")

	if err := cmd.Err(); err != nil {
		return nil, err
	}

	keys := lo.Map(cmd.Val(), func(key string, _ int) string {
		sf, _ := discord.ParseSnowflake(key)

		return c.key(discord.ChannelID(sf))
	})

	return MGet[discord.Channel](c.re, keys...)
}

func (c *Channel) PrivateChannels() ([]discord.Channel, error) {
	keys, err := c.re.Keys("channels:private:*")

	if err != nil {
		return nil, err
	}

	return MGet[discord.Channel](c.re, keys...)
}

func (c *Channel) ChannelSet(channel *discord.Channel, update bool) error {
	key := c.key(channel.ID)

	exists, err := c.re.Exists(key)

	if err != nil {
		return err
	}

	if exists && !update {
		return nil
	}

	if err := c.re.Set(key, *channel); err != nil {
		return err
	}

	switch channel.Type {
	case discord.DirectMessage:
		// Safety bound check.
		if len(channel.DMRecipients) != 1 {
			return fmt.Errorf("DirectMessage channel %d doesn't have 1 recipient", channel.ID)
		}

		return c.re.Set("channels:private:"+channel.DMRecipients[0].ID.String(), channel.ID)
	case discord.GroupDM:
		return c.re.client.SAdd(context.Background(), "guilds:0:channels", channel.ID.String()).Err()
	}

	return c.re.client.SAdd(context.Background(), "guilds:"+channel.GuildID.String()+":channels", channel.ID.String()).Err()
}

func (c *Channel) ChannelRemove(channel *discord.Channel) error {
	if err := c.re.Delete(c.key(channel.ID)); err != nil {
		return err
	}

	switch channel.Type {
	case discord.DirectMessage:
		if len(channel.DMRecipients) != 1 {
			return fmt.Errorf("DirectMessage channel %d doesn't have 1 recipient", channel.ID)
		}

		return c.re.Delete("channels:private:" + channel.DMRecipients[0].ID.String())
	case discord.GroupDM:
		// Set guild channels for 0 to remove group dm
		return c.re.client.SRem(context.Background(), "guilds:0:channels", channel.ID.String()).Err()
	}

	return c.re.client.SRem(context.Background(), "guilds:"+channel.GuildID.String()+":channels", channel.ID.String()).Err()
}
