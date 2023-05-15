package redis

import (
	"context"
	"encoding/json"
	"github.com/diamondburned/arikawa/v3/discord"
)

type Emoji struct {
	re *Redis
}

func (e *Emoji) Reset() error {
	return e.re.PatternDel(makeKey("emojis", discord.NullGuildID))
}

func (e *Emoji) key(guildID discord.GuildID, emojiID discord.EmojiID) string {
	return makeKey("emojis", guildID, emojiID)
}

func (e *Emoji) Emoji(guildID discord.GuildID, emojiID discord.EmojiID) (*discord.Emoji, error) {
	return Get[*discord.Emoji](e.re, e.key(guildID, emojiID))
}

func (e *Emoji) Emojis(guildID discord.GuildID) ([]discord.Emoji, error) {
	keys, err := e.re.Keys(e.key(guildID, discord.NullEmojiID))

	if err != nil {
		return nil, err
	}

	return MGet[discord.Emoji](e.re, keys...)
}

func (e *Emoji) EmojiSet(guildID discord.GuildID, emojis []discord.Emoji, update bool) error {
	// TODO: Check emojis to see if we already have ANY for the guild?

	m := make(map[string]interface{})

	for _, emoji := range emojis {
		val, err := json.Marshal(emoji)

		if err != nil {
			return err
		}

		m[e.key(guildID, emoji.ID)] = val
	}

	return e.re.client.MSet(context.Background(), m).Err()
}
