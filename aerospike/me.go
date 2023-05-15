package aerospike

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/diamondburned/arikawa/v3/discord"
)

type Me struct {
	as *Aerospike
}

func (m *Me) key() (*aerospike.Key, error) {
	return m.as.newKey("users", "me")
}

func (m *Me) Me() (*discord.User, error) {
	k, err := m.key()

	if err != nil {
		return nil, err
	}

	return Get[*discord.User](m.as, k)
}

func (m *Me) MyselfSet(me discord.User, update bool) error {
	return m.as.Set("users", "me", me)
}

func (m *Me) Reset() error {
	return m.as.Delete("users", "me")
}
