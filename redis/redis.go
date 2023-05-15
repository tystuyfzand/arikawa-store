package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/diamondburned/arikawa/v3/state/store"
	"github.com/redis/go-redis/v9"
)

var patternDel = redis.NewScript(`
for _, k in ipairs(redis.call('KEYS', ARGV[1])) do
	redis.call('DEL', k)
end

return 1
`)

func New(client *redis.Client) *store.Cabinet {
	re := &Redis{
		client: client,
	}

	return &store.Cabinet{
		MeStore:         &Me{re: re},
		ChannelStore:    &Channel{re: re},
		EmojiStore:      &Emoji{re: re},
		GuildStore:      &Guild{re: re},
		MemberStore:     &Member{re: re},
		MessageStore:    &Message{re: re},
		PresenceStore:   &Presence{re: re},
		RoleStore:       &Role{re: re},
		VoiceStateStore: &VoiceState{re: re},
	}
}

type Redis struct {
	client *redis.Client
}

// Keys is a shorthand for redis.Keys which returns the values as multiple params.
func (r *Redis) Keys(pattern string) ([]string, error) {
	keys := r.client.Keys(context.Background(), pattern)

	if err := keys.Err(); err != nil {
		return nil, err
	}

	return keys.Val(), nil
}

// Exists is a shorthand for redis.Exists which returns the values as params.
func (r *Redis) Exists(key string) (bool, error) {
	cmd := r.client.Exists(context.Background(), key)

	if err := cmd.Err(); err != nil {
		return false, err
	}

	return cmd.Val() == 1, nil
}

// Set is a wrapper for redis.Set which uses json to marshal objects
func (r *Redis) Set(key string, val any) error {
	b, err := json.Marshal(val)

	if err != nil {
		return err
	}

	cmd := r.client.Set(context.Background(), key, b, 0)

	return cmd.Err()
}

// Delete wraps redis.Del and returns the error
func (r *Redis) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func (r *Redis) PatternDel(pattern string) error {
	return patternDel.Run(context.Background(), r.client, []string{}, pattern).Err()
}

func Get[V any](re *Redis, key string) (val V, err error) {
	cmd := re.client.Get(context.Background(), key)

	if err = cmd.Err(); err != nil {
		return
	}

	var b []byte

	if b, err = cmd.Bytes(); err != nil {
		return
	}

	err = json.Unmarshal(b, &val)

	return
}

func MGet[V any](re *Redis, keys ...string) (val []V, err error) {
	cmd := re.client.MGet(context.Background(), keys...)

	if err = cmd.Err(); err != nil {
		return
	}

	for _, v := range cmd.Val() {
		data, err := Unmarshal[V](v)

		if err != nil {
			return nil, err
		}

		val = append(val, data)
	}

	return
}

func Unmarshal[V any](data any) (val V, err error) {
	var bytes []byte

	if b, ok := data.([]byte); ok {
		bytes = b
	} else if s, ok := data.(string); ok {
		bytes = []byte(s)
	} else {
		err = fmt.Errorf("unexpected type %T", data)
		return
	}

	err = json.Unmarshal(bytes, &val)

	return
}

type Snowflake interface {
	IsValid() bool
	String() string
}

func makeKey(prefix string, ids ...Snowflake) string {
	str := prefix

	for _, id := range ids {
		str += ":"

		if id.IsValid() {
			str += id.String()
		} else {
			str += "*"
		}
	}

	return str
}
