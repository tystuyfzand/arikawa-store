package aerospike

import "github.com/diamondburned/arikawa/v3/discord"

type Message struct {
	as          *Aerospike
	maxMessages int
}

func (m *Message) Reset() error {
	return m.as.Truncate("messages")
}

func (m *Message) MaxMessages() int {
	return m.maxMessages
}

func (m *Message) Message(channelID discord.ChannelID, messageID discord.MessageID) (*discord.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Message) Messages(channelID discord.ChannelID) ([]discord.Message, error) {
	// Retrieve messages, deserialize, list
	return nil, nil
}

func (m *Message) MessageSet(message *discord.Message, update bool) error {
	//TODO implement me
	panic("implement me")
	// Append to list?
}

func (m *Message) MessageRemove(channelID discord.ChannelID, messageID discord.MessageID) error {
	//TODO implement me
	panic("implement me")
}
