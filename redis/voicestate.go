package redis

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

type VoiceState struct {
	re *Redis
}

func (v *VoiceState) Reset() error {
	return v.re.PatternDel(makeKey("voicestates", discord.NullGuildID))
}

func (v *VoiceState) key(guildID discord.GuildID, userID discord.UserID) string {
	return makeKey("voicestates", guildID, userID)
}

func (v *VoiceState) VoiceState(guildID discord.GuildID, userID discord.UserID) (*discord.VoiceState, error) {
	return Get[*discord.VoiceState](v.re, v.key(guildID, userID))
}

func (v *VoiceState) VoiceStates(guildID discord.GuildID) ([]discord.VoiceState, error) {
	keys, err := v.re.Keys(v.key(guildID, discord.NullUserID))

	if err != nil {
		return nil, err
	}

	return MGet[discord.VoiceState](v.re, keys...)
}

func (v *VoiceState) VoiceStateSet(guildID discord.GuildID, state *discord.VoiceState, update bool) error {
	key := v.key(guildID, state.UserID)

	exists, err := v.re.Exists(key)

	if err != nil {
		return err
	}

	if exists && !update {
		return nil
	}

	return v.re.Set(key, *state)
}

func (v *VoiceState) VoiceStateRemove(guildID discord.GuildID, userID discord.UserID) error {
	return v.re.Delete(v.key(guildID, userID))
}
