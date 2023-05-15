package aerospike

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/diamondburned/arikawa/v3/discord"
)

type Emoji struct {
	as *Aerospike
}

func (e *Emoji) Reset() error {
	return e.as.Truncate("emojis")
}

func (e *Emoji) Emoji(guildID discord.GuildID, emojiID discord.EmojiID) (*discord.Emoji, error) {
	k, err := e.as.newKey("emojis", guildID)

	if err != nil {
		return nil, err
	}

	return Get[*discord.Emoji](e.as, k, emojiID.String())
}

func (e *Emoji) Emojis(guildID discord.GuildID) ([]discord.Emoji, error) {
	return GetKeyBins[discord.Emoji](e.as, "emojis", guildID)
}

func (e *Emoji) EmojiSet(guildID discord.GuildID, emojis []discord.Emoji, update bool) error {
	k, err := e.as.newKey("emojis", guildID)

	if err != nil {
		return err
	}

	exists, err := e.as.client.Exists(nil, k)

	if err != nil {
		return err
	}

	if exists && !update {
		return nil
	}

	bins := make(aerospike.BinMap)

	for _, emoji := range emojis {
		bins[emoji.ID.String()], err = MarshalBin(emoji)

		if err != nil {
			return err
		}
	}

	return e.as.client.Add(nil, k, bins)
}
