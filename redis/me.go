package redis

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

type Me struct {
	re *Redis
}

func (m *Me) Reset() error {
	return m.re.Delete("me")
}

func (m *Me) Me() (*discord.User, error) {
	return Get[*discord.User](m.re, "me")
}

func (m *Me) MyselfSet(me discord.User, update bool) error {
	exists, err := m.re.Exists("me")

	if err != nil {
		return err
	}

	if exists && !update {
		return nil
	}

	return m.re.Set("me", me)
}
