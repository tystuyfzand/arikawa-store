package aerospike

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

type Role struct {
	as *Aerospike
}

func (r *Role) Reset() error {
	return r.as.Truncate("roles")
}

func (r *Role) Role(guildID discord.GuildID, roleID discord.RoleID) (*discord.Role, error) {
	k, err := r.as.newKey("roles", guildID)

	if err != nil {
		return nil, err
	}

	return Get[*discord.Role](r.as, k, roleID.String())
}

func (r *Role) Roles(guildID discord.GuildID) ([]discord.Role, error) {
	return GetKeyBins[discord.Role](r.as, "roles", guildID)
}

func (r *Role) RoleSet(guildID discord.GuildID, role *discord.Role, update bool) error {
	return r.as.SetBin("roles", role.ID.String(), guildID, role)
}

func (r *Role) RoleRemove(guildID discord.GuildID, roleID discord.RoleID) error {
	return r.as.SetBin("roles", roleID.String(), guildID, nil)
}
