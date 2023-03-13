package types

import (
	"fmt"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId:          PortID,
		LotteryPotsList: []LotteryPots{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}
	// Check for duplicated ID in lotteryPots
	lotteryPotsIdMap := make(map[uint64]bool)
	lotteryPotsCount := gs.GetLotteryPotsCount()
	for _, elem := range gs.LotteryPotsList {
		if _, ok := lotteryPotsIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for lotteryPots")
		}
		if elem.Id >= lotteryPotsCount {
			return fmt.Errorf("lotteryPots id should be lower or equal than the last id")
		}
		lotteryPotsIdMap[elem.Id] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
