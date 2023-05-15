package redis

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

type Guild struct {
	re *Redis
}

func (g *Guild) Reset() error {
	return g.re.PatternDel("guilds:*")
}

func (g *Guild) key(guildID discord.GuildID) string {
	return makeKey("guilds", guildID)
}

func (g *Guild) Guild(guildID discord.GuildID) (*discord.Guild, error) {
	return Get[*discord.Guild](g.re, g.key(guildID))
}

func (g *Guild) Guilds() ([]discord.Guild, error) {
	keys, err := g.re.Keys(g.key(discord.NullGuildID))

	if err != nil {
		return nil, err
	}

	return MGet[discord.Guild](g.re, keys...)
}

func (g *Guild) GuildSet(guild *discord.Guild, update bool) error {
	key := g.key(guild.ID)

	exists, err := g.re.Exists(key)

	if err != nil {
		return err
	}

	if exists && !update {
		return nil
	}

	return g.re.Set(key, *guild)
}

func (g *Guild) GuildRemove(guildID discord.GuildID) error {
	return g.re.Delete(g.key(guildID))
}
