package debug

import (
	"fmt"

	"github.com/s7techlab/cckit/state"
)

// DeleteStateByPrefixes deletes from state entries with matching key prefix
// raw function, do not use State wrappers, like encryption
func DeleteStateByPrefixes(s state.State, prefixes []string) (map[string]uint32, error) {
	prefixMatches := make(map[string]uint32)

	for _, prefix := range prefixes {

		keys, err := s.Keys(prefix)
		if err != nil {
			return nil, fmt.Errorf(`keys: %w`, err)
		}

		prefixMatches[prefix] = 0
		for _, key := range keys {
			if err = s.Delete(key); err != nil {
				return nil, err
			}
			prefixMatches[prefix]++
		}
	}

	return prefixMatches, nil
}
