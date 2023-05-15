package aerospike

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

type Member struct {
	as *Aerospike
}

func (m *Member) Reset() error {
	return m.as.Truncate("members")
}

func (m *Member) Member(guildID discord.GuildID, userID discord.UserID) (*discord.Member, error) {
	k, err := m.as.newKey("members", guildID)

	if err != nil {
		return nil, err
	}

	return Get[*discord.Member](m.as, k, userID.String())
}

func (m *Member) Members(guildID discord.GuildID) ([]discord.Member, error) {
	return GetKeyBins[discord.Member](m.as, "members", guildID)
}

func (m *Member) MemberSet(guildID discord.GuildID, member *discord.Member, update bool) error {
	return m.as.SetBin("members", member.User.ID.String(), guildID, member)
}

func (m *Member) MemberRemove(guildID discord.GuildID, userID discord.UserID) error {
	return m.as.SetBin("members", userID.String(), guildID, nil)
}
