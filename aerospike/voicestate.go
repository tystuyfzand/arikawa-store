package aerospike

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

type VoiceState struct {
	as *Aerospike
}

func (v *VoiceState) Reset() error {
	return v.as.Truncate("voicestates")
}

func (v *VoiceState) VoiceState(guildID discord.GuildID, userID discord.UserID) (*discord.VoiceState, error) {
	k, err := v.as.newKey("voicestates", guildID)

	if err != nil {
		return nil, err
	}

	return Get[*discord.VoiceState](v.as, k, userID.String())
}

func (v *VoiceState) VoiceStates(guildID discord.GuildID) ([]discord.VoiceState, error) {
	return GetKeyBins[discord.VoiceState](v.as, "voicestates", guildID)
}

func (v *VoiceState) VoiceStateSet(guildID discord.GuildID, state *discord.VoiceState, update bool) error {
	return v.as.SetBin("voicestates", state.UserID.String(), guildID, state)
}

func (v *VoiceState) VoiceStateRemove(guildID discord.GuildID, userID discord.UserID) error {
	return v.as.SetBin("voicestates", userID.String(), guildID, nil)
}
