package aerospike

import (
	"encoding/json"
	"errors"
	aero "github.com/aerospike/aerospike-client-go"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state/store"
	"reflect"
)

type Aerospike struct {
	client    *aero.Client
	namespace string
}

func (a *Aerospike) newKey(set string, key any) (*aero.Key, error) {
	if uintVal, ok := castToUint64(key); ok {
		key = uintVal
	}

	return aero.NewKey(a.namespace, set, key)
}

func New(client *aero.Client) *store.Cabinet {
	as := &Aerospike{
		client: client,
	}

	return &store.Cabinet{
		MeStore:      &Me{as: as},
		ChannelStore: &Channel{as: as},
		EmojiStore:   &Emoji{as: as},
		GuildStore:   &Guild{as: as},
		MemberStore:  &Member{as: as},
		// MessageStore: &Message{as: as},
		PresenceStore:   &Presence{as: as},
		RoleStore:       &Role{as: as},
		VoiceStateStore: &VoiceState{as: as},
	}
}

func (a *Aerospike) listPolicy() *aero.ListPolicy {
	return aero.NewListPolicy(aero.ListOrderUnordered, aero.ListWriteFlagsAddUnique|aero.ListWriteFlagsNoFail)
}

func (a *Aerospike) ListAdd(set, binName string, key, v any) error {
	asKey, err := a.newKey(set, key)

	if err != nil {
		return err
	}

	if uintVal, ok := castToUint64(v); ok {
		v = uintVal
	}

	_, err = a.client.Operate(nil, asKey, aero.ListAppendWithPolicyOp(a.listPolicy(), binName, v))

	return err
}

func (a *Aerospike) ListRemove(set, binName string, key, v any) error {
	asKey, err := a.newKey(set, key)

	if err != nil {
		return err
	}

	if uintVal, ok := castToUint64(v); ok {
		v = uintVal
	}

	_, err = a.client.Operate(nil, asKey, aero.ListRemoveByValueOp(binName, v, aero.ListReturnTypeNone))

	return err
}

func (a *Aerospike) Truncate(set string) error {
	return a.client.Truncate(nil, a.namespace, set, nil)
}

func (a *Aerospike) Set(set string, key, v any) error {
	return a.SetBin(set, "value", key, v)
}

func (a *Aerospike) SetBin(set, bin string, key, v any) error {
	asKey, err := a.newKey(set, key)

	if err != nil {
		return err
	}

	if uintVal, ok := castToUint64(v); ok {
		v = uintVal
	} else if v != nil {
		var err error

		v, err = json.Marshal(v)

		if err != nil {
			return err
		}
	}

	policy := aero.NewWritePolicy(0, 0)
	policy.SendKey = true

	return a.client.AddBins(policy, asKey, aero.NewBin(bin, v))
}

func (a *Aerospike) Exists(set string, k any) (bool, error) {
	asKey, err := a.newKey(set, k)

	if err != nil {
		return false, err
	}

	return a.client.Exists(nil, asKey)
}

func castToUint64(v any) (uint64, bool) {
	switch val := v.(type) {
	case discord.Snowflake:
		v = uint64(val)
	case discord.AppID:
		v = uint64(val)
	case discord.AttachmentID:
		v = uint64(val)
	case discord.AuditLogEntryID:
		v = uint64(val)
	case discord.ChannelID:
		v = uint64(val)
	case discord.CommandID:
		v = uint64(val)
	case discord.EmojiID:
		v = uint64(val)
	case discord.GuildID:
		v = uint64(val)
	case discord.IntegrationID:
		v = uint64(val)
	case discord.InteractionID:
		v = uint64(val)
	case discord.MessageID:
		v = uint64(val)
	case discord.RoleID:
		v = uint64(val)
	case discord.StageID:
		v = uint64(val)
	case discord.StickerID:
		v = uint64(val)
	case discord.StickerPackID:
		v = uint64(val)
	case discord.TagID:
		v = uint64(val)
	case discord.TeamID:
		v = uint64(val)
	case discord.UserID:
		v = uint64(val)
	case discord.WebhookID:
		v = uint64(val)
	case discord.EventID:
		v = uint64(val)
	case discord.EntityID:
		v = uint64(val)
	}

	if val, ok := v.(uint64); ok {
		return val, ok
	}

	return 0, false
}

func (a *Aerospike) Delete(set string, key any) error {
	asKey, err := a.newKey(set, key)

	if err != nil {
		return err
	}

	_, err = a.client.Delete(nil, asKey)

	return err
}

var (
	snowflakeType = reflect.TypeOf((*discord.Snowflake)(nil)).Elem()
)

func Get[V any](as *Aerospike, key *aero.Key, binNames ...string) (val V, err error) {
	var record *aero.Record

	record, err = as.client.Get(nil, key, binNames...)

	if err != nil {
		return
	}

	if binNames == nil {
		binNames = []string{"value"}
	}

	b := record.Bins[binNames[0]]

	return UnmarshalBin[V](b)
}

func MarshalBin(v any) (any, error) {
	if uintVal, ok := castToUint64(v); ok {
		return uintVal, nil
	} else if str, ok := v.(string); ok {
		return str, nil
	}

	return json.Marshal(v)
}

func UnmarshalBin[V any](b any) (val V, err error) {
	vType := reflect.TypeOf(val)

	switch bType := b.(type) {
	case []byte:
		err = json.Unmarshal(bType, &val)
		return
	case []uint64: // Only used for snowflake slices
		if vType.Kind() == reflect.Slice {
			// Check snowflakes
			vType = vType.Elem()

			// Direct conversion
			if vType.Kind() == reflect.Uint64 {
				val = b.(V)
				return
			}

			// Snowflake types
			if !vType.ConvertibleTo(snowflakeType) {
				err = errors.New("cannot convert []int64 to anything other than snowflake")
				return
			}

			refVal := reflect.ValueOf(val)

			for _, val := range bType {
				refVal = reflect.Append(refVal, reflect.ValueOf(val).Convert(vType))
			}

			val = refVal.Interface().(V)
		}
	default:
		err = errors.New("unknown type")
	}

	return
}

func firstErr(err []error) error {
	for _, e := range err {
		if e != nil {
			return e
		}
	}

	return nil
}
