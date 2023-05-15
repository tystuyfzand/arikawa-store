package aerospike

import (
	"github.com/diamondburned/arikawa/v3/state/store"
)

func GetKeyBins[V any](as *Aerospike, set string, key any) ([]V, error) {
	k, err := as.newKey(set, key)

	if err != nil {
		return nil, err
	}

	record, err := as.client.Get(nil, k)

	if err != nil {
		return nil, store.ErrNotFound
	}

	val := make([]V, len(record.Bins))
	i := 0

	for _, data := range record.Bins {
		val[i], err = UnmarshalBin[V](data)

		if err != nil {
			return nil, err
		}

		i++
	}

	return val, nil
}
