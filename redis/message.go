package redis

import (
	"context"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

var messageCleanupScript = redis.NewScript(`
	-- KEYS[1] the sorted set
	-- ARGV[1] the max set size
	-- returns number of elements removed
	--
	local setSize = redis.call('ZCARD', KEYS[1])
	local maxSize = tonumber(ARGV[1])
	if setSize > maxSize then
		local pop = redis.call('ZPOPMIN', KEYS[1], setSize - maxSize)

		for i = 1, #pop, 2 do
			redis.call('DEL', string.sub(KEYS[1], 0, -5) .. pop[i])
		end

		return #pop / 2
	else
		return 0    
	end
`)

// Message uses Redis Sorted Sets.
// Messages are stored in messages:CHANNELID:MESSAGEID, and a reference kept in
// messages:CHANNELID:list with associated unix timestamps (float64 is the max value)
// When the collection gets too big, an item is automatically pulled off the start and removed.
// See messageCleanupScript above for how we do it.
type Message struct {
	re          *Redis
	maxMessages int
}

func (m *Message) Reset() error {
	return m.re.PatternDel("messages:*")
}

func (m *Message) MaxMessages() int {
	return m.maxMessages
}

func (m *Message) key(channelID discord.ChannelID, messageID discord.MessageID) string {
	return makeKey("messages", channelID, messageID)
}

func (m *Message) Message(channelID discord.ChannelID, messageID discord.MessageID) (*discord.Message, error) {
	return Get[*discord.Message](m.re, m.key(channelID, messageID))
}

func (m *Message) Messages(channelID discord.ChannelID) ([]discord.Message, error) {
	ids := m.re.client.ZRange(context.Background(), "messages:"+channelID.String()+":list", 0, int64(m.maxMessages))

	if err := ids.Err(); err != nil {
		return nil, err
	}

	keys := lo.Map(ids.Val(), func(item string, _ int) string {
		return "messages:" + channelID.String() + ":" + item
	})

	return MGet[discord.Message](m.re, keys...)
}

func (m *Message) MessageSet(message *discord.Message, update bool) error {
	key := m.key(message.ChannelID, message.ID)

	exists, err := m.re.Exists(key)

	if err != nil {
		return err
	}

	if exists && !update {
		return nil
	}

	err = m.re.Set(key, *message)

	if err != nil {
		return err
	}

	listKey := "messages:" + message.ChannelID.String() + ":list"

	cmd := m.re.client.ZAdd(context.Background(), listKey, redis.Z{
		Member: message.ID.String(),
		Score:  float64(message.Timestamp.Time().Unix()),
	})

	if err = cmd.Err(); err != nil {
		return err
	}

	scriptCmd := messageCleanupScript.Run(context.Background(), m.re.client, []string{listKey}, m.maxMessages)

	if err := scriptCmd.Err(); err != nil {
		return err
	}

	return nil
}

func (m *Message) MessageRemove(channelID discord.ChannelID, messageID discord.MessageID) error {
	key := m.key(channelID, messageID)

	err := m.re.Delete(key)

	if err != nil {
		return err
	}

	return m.re.client.ZRem(context.Background(), "messages:"+channelID.String()+":list", messageID.String()).Err()
}
