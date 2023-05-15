package redis

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

type Presence struct {
	re *Redis
}

func (p *Presence) Reset() error {
	return p.re.PatternDel(makeKey("presence", discord.NullGuildID))
}

func (p *Presence) key(guildID discord.GuildID, userID discord.UserID) string {
	return makeKey("presence", guildID, userID)
}

func (p *Presence) Presence(guildID discord.GuildID, userID discord.UserID) (*discord.Presence, error) {
	return Get[*discord.Presence](p.re, p.key(guildID, userID))
}

func (p *Presence) Presences(guildID discord.GuildID) ([]discord.Presence, error) {
	keys, err := p.re.Keys(p.key(guildID, discord.NullUserID))

	if err != nil {
		return nil, err
	}

	return MGet[discord.Presence](p.re, keys...)
}

func (p *Presence) PresenceSet(guildID discord.GuildID, presence *discord.Presence, update bool) error {
	key := p.key(guildID, presence.User.ID)

	exists, err := p.re.Exists(key)

	if err != nil {
		return err
	}

	if exists && !update {
		return nil
	}

	return p.re.Set(key, *presence)
}

func (p *Presence) PresenceRemove(guildID discord.GuildID, userID discord.UserID) error {
	return p.re.Delete(p.key(guildID, userID))
}
