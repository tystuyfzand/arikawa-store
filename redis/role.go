package redis

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

type Role struct {
	re *Redis
}

func (r *Role) Reset() error {
	return r.re.PatternDel(makeKey("roles", discord.NullRoleID))
}

func (r *Role) key(guildID discord.GuildID, roleID discord.RoleID) string {
	return makeKey("roles", guildID, roleID)
}

func (r *Role) Role(guildID discord.GuildID, roleID discord.RoleID) (*discord.Role, error) {
	return Get[*discord.Role](r.re, r.key(guildID, roleID))
}

func (r *Role) Roles(guildID discord.GuildID) ([]discord.Role, error) {
	keys, err := r.re.Keys(r.key(guildID, discord.NullRoleID))

	if err != nil {
		return nil, err
	}

	return MGet[discord.Role](r.re, keys...)
}

func (r *Role) RoleSet(guildID discord.GuildID, role *discord.Role, update bool) error {
	key := r.key(guildID, role.ID)

	exists, err := r.re.Exists(key)

	if err != nil {
		return err
	}

	if exists && !update {
		return nil
	}

	return r.re.Set(key, *role)
}

func (r *Role) RoleRemove(guildID discord.GuildID, roleID discord.RoleID) error {
	return r.re.Delete(r.key(guildID, roleID))
}
