package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"lottery/x/lottery/keeper"
	"lottery/x/lottery/types"
)

func SimulateMsgCreateLottery(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgCreateLottery{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the CreateLottery simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "CreateLottery simulation not implemented"), nil, nil
	}
}
