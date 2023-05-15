package aerospike

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

type Presence struct {
	as *Aerospike
}

func (p *Presence) Reset() error {
	return p.as.Truncate("presences")
}

func (p *Presence) Presence(guildID discord.GuildID, userID discord.UserID) (*discord.Presence, error) {
	k, err := p.as.newKey("presence", guildID)

	if err != nil {
		return nil, err
	}

	return Get[*discord.Presence](p.as, k, userID.String())
}

func (p *Presence) Presences(guildID discord.GuildID) ([]discord.Presence, error) {
	return GetKeyBins[discord.Presence](p.as, "presences", guildID)
}

func (p *Presence) PresenceSet(guildID discord.GuildID, presence *discord.Presence, update bool) error {
	return p.as.SetBin("presences", presence.User.ID.String(), guildID, presence)
}

func (p *Presence) PresenceRemove(guildID discord.GuildID, userID discord.UserID) error {
	return p.as.SetBin("presences", userID.String(), guildID, nil)
}
