package redis

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

type Member struct {
	re *Redis
}

func (m *Member) Reset() error {
	return m.re.PatternDel(makeKey("members", discord.NullUserID))
}

func (m *Member) key(guildID discord.GuildID, userID discord.UserID) string {
	return makeKey("members", guildID, userID)
}

func (m *Member) Member(guildID discord.GuildID, userID discord.UserID) (*discord.Member, error) {
	return Get[*discord.Member](m.re, m.key(guildID, userID))
}

func (m *Member) Members(guildID discord.GuildID) ([]discord.Member, error) {
	keys, err := m.re.Keys(m.key(guildID, discord.NullUserID))

	if err != nil {
		return nil, err
	}

	return MGet[discord.Member](m.re, keys...)
}

func (m *Member) MemberSet(guildID discord.GuildID, member *discord.Member, update bool) error {
	key := m.key(guildID, member.User.ID)

	exists, err := m.re.Exists(key)

	if err != nil {
		return nil
	}

	if exists && !update {
		return nil
	}

	return m.re.Set(key, *member)
}

func (m *Member) MemberRemove(guildID discord.GuildID, userID discord.UserID) error {
	return m.re.Delete(m.key(guildID, userID))
}
