package lottery

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"lottery/x/lottery/keeper"
	"lottery/x/lottery/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the lotteryPots
	for _, elem := range genState.LotteryPotsList {
		k.SetLotteryPots(ctx, elem)
	}

	// Set lotteryPots count
	k.SetLotteryPotsCount(ctx, genState.LotteryPotsCount)
	// this line is used by starport scaffolding # genesis/module/init
	k.SetPort(ctx, genState.PortId)
	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if !k.IsBound(ctx, genState.PortId) {
		// module binds to the port on InitChain
		// and claims the returned capability
		err := k.BindPort(ctx, genState.PortId)
		if err != nil {
			panic("could not claim port capability: " + err.Error())
		}
	}
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PortId = k.GetPort(ctx)
	genesis.LotteryPotsList = k.GetAllLotteryPots(ctx)
	genesis.LotteryPotsCount = k.GetLotteryPotsCount(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
