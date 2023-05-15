package aerospike

import (
	"errors"
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state/store"
)

type Channel struct {
	as *Aerospike
}

func (c *Channel) Reset() error {
	errs := []error{
		c.as.Truncate("channels"),
		c.as.Truncate("privates"),
	}

	// When guilds are reset, so are channels on those guilds.
	// TODO: Is this necessary to reset if the whole cabinet is?

	return firstErr(errs)
}

func (c *Channel) Channel(id discord.ChannelID) (*discord.Channel, error) {
	k, err := c.as.newKey("channels", id)

	if err != nil {
		return nil, err
	}

	return Get[*discord.Channel](c.as, k)
}

func (c *Channel) CreatePrivateChannel(recipient discord.UserID) (*discord.Channel, error) {
	k, err := c.as.newKey("privates", recipient)

	if err != nil {
		return nil, err
	}

	privateID, err := Get[discord.ChannelID](c.as, k)

	if err != nil {
		return nil, store.ErrNotFound
	}

	return c.Channel(privateID)
}

func (c *Channel) Channels(id discord.GuildID) ([]discord.Channel, error) {
	k, err := c.as.newKey("guilds", id)

	if err != nil {
		return nil, err
	}

	channelIDs, err := Get[[]discord.ChannelID](c.as, k, "channels")

	if err != nil {
		return nil, store.ErrNotFound
	}

	channels := make([]discord.Channel, len(channelIDs))

	for i, channelID := range channelIDs {
		channel, err := c.Channel(channelID)

		if err != nil {
			return nil, err
		}

		channels[i] = *channel
	}

	return channels, nil
}

func (c *Channel) PrivateChannels() ([]discord.Channel, error) {
	scan, err := c.as.client.ScanAll(nil, c.as.namespace, "privates", "value")

	if err != nil {
		return nil, err
	}

	defer scan.Close()

	channels := make([]discord.Channel, 0)

L:
	for {
		select {
		case rec := <-scan.Records:
			if rec == nil {
				break L
			}

			channel, err := UnmarshalBin[discord.Channel](rec.Bins["value"].([]byte))

			if err != nil {
				return nil, err
			}

			channels = append(channels, channel)
		case err := <-scan.Errors:
			return nil, err
		}
	}

	return channels, nil
}

func (c *Channel) listPolicy() *aerospike.ListPolicy {
	return aerospike.NewListPolicy(aerospike.ListOrderUnordered, aerospike.ListWriteFlagsAddUnique|aerospike.ListWriteFlagsNoFail)
}

func (c *Channel) ChannelSet(channel *discord.Channel, update bool) error {
	if err := c.as.Set("channels", channel.ID, *channel); err != nil {
		return err
	}

	switch channel.Type {
	case discord.DirectMessage:
		// Safety bound check.
		if len(channel.DMRecipients) != 1 {
			return fmt.Errorf("DirectMessage channel %d doesn't have 1 recipient", channel.ID)
		}

		return c.as.Set("privates", channel.DMRecipients[0].ID, channel.ID)
	case discord.GroupDM:
		return c.as.ListAdd("guilds", "channels", uint64(0), channel.ID)
	}

	if !channel.GuildID.IsValid() {
		return errors.New("invalid guildID for guild channel")
	}

	return c.as.ListAdd("guilds", "channels", channel.GuildID, channel.ID)
}

func (c *Channel) ChannelRemove(channel *discord.Channel) error {
	err := c.as.Delete("channels", channel.ID)

	if err != nil {
		return err
	}

	switch channel.Type {
	case discord.DirectMessage:
		if len(channel.DMRecipients) != 1 {
			return fmt.Errorf("DirectMessage channel %d doesn't have 1 recipient", channel.ID)
		}

		return c.as.Delete("privates", channel.DMRecipients[0].ID)
	case discord.GroupDM:
		// Set guild channels for 0 to remove group dm
		return c.as.ListRemove("guilds", "channels", uint64(0), channel.ID)
	}

	return c.as.ListRemove("guilds", "channels", channel.GuildID, channel.ID)
}
