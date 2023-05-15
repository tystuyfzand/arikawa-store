package aerospike

import "github.com/diamondburned/arikawa/v3/discord"

type Guild struct {
	as *Aerospike
}

func (g *Guild) Reset() error {
	return g.as.Truncate("guilds")
}

func (g *Guild) Guild(guildID discord.GuildID) (*discord.Guild, error) {
	key, err := g.as.newKey("guilds", guildID)

	if err != nil {
		return nil, err
	}

	return Get[*discord.Guild](g.as, key)
}

func (g *Guild) Guilds() ([]discord.Guild, error) {
	scan, err := g.as.client.ScanAll(nil, g.as.namespace, "guilds", "value")

	if err != nil {
		return nil, err
	}

	defer scan.Close()

	guilds := make([]discord.Guild, 0)

L:
	for {
		select {
		case rec := <-scan.Records:
			if rec == nil {
				break L
			}

			guild, err := UnmarshalBin[discord.Guild](rec.Bins["value"].([]byte))

			if err != nil {
				return nil, err
			}

			guilds = append(guilds, guild)
		case err := <-scan.Errors:
			return nil, err
		}
	}

	return guilds, nil
}

func (g *Guild) GuildSet(guild *discord.Guild, update bool) error {
	exists, err := g.as.Exists("guilds", guild.ID)

	if err != nil {
		return err
	}

	if exists && !update {
		return nil
	}

	return g.as.Set("guilds", guild.ID, guild)
}

func (g *Guild) GuildRemove(guildID discord.GuildID) error {
	return g.as.Delete("guilds", guildID)
}
