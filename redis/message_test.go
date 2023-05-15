package redis

import (
	"context"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestMessages(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr: ":6379",
		DB:   0,
	})

	if err := c.Ping(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}

	messages := &Message{
		re:          &Redis{client: c},
		maxMessages: 100,
	}

	channelID := discord.ChannelID(123)

	offset := 0

	addMessages := func(max int) {
		for i := 0; i < max; i++ {
			err := messages.MessageSet(&discord.Message{
				ID:        discord.MessageID(i + offset + 1),
				ChannelID: channelID,
				Timestamp: discord.NewTimestamp(time.Now().Add(time.Duration(i+offset*2) * time.Second)),
			}, true)

			if err != nil {
				t.Fatal("Unable to set messages", err)
			}
		}

		offset += max
	}

	addMessages(100)

	t.Log("Added 100 messages")

	checkLength := func() {
		t.Log("Checking message size")
		m, err := messages.Messages(channelID)

		if err != nil {
			t.Fatal("Unable to retrieve", err)
		}

		if len(m) > messages.maxMessages {
			t.Fatal("Invalid length!", len(m))
		}
	}
	checkLength()

	time.Sleep(1 * time.Second)

	t.Log("Adding 5")

	addMessages(5)

	checkLength()

	t.Log("Removing message")

	err := messages.MessageRemove(channelID, discord.MessageID(offset-1))

	if err != nil {
		t.Fatal("Unable to remove message", err)
	}

	err = messages.Reset()

	if err != nil {
		t.Fatal("Unable to cleanup", err)
	}
}
