package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ScavengeList: []Scavenge{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in scavenge
	scavengeIndexMap := make(map[string]struct{})

	for _, elem := range gs.ScavengeList {
		index := string(ScavengeKey(elem.Index))
		if _, ok := scavengeIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for scavenge")
		}
		scavengeIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
