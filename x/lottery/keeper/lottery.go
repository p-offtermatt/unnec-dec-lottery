package keeper

import (
	"encoding/binary"

	"lottery/x/lottery/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetLotteryCount get the total number of lotteries
func (k Keeper) GetLotteryCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.LotteryCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetLotteryCount set the total number of Lotteries
func (k Keeper) SetLotteryCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.LotteryCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendLottery appends a lottery in the store with a new id and update the count
func (k Keeper) AppendLottery(
	ctx sdk.Context,
	lottery types.Lottery,
) uint64 {
	count := k.GetLotteryCount(ctx)

	// Set the ID of the appended value
	lottery.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LotteryKey))
	appendedValue := k.cdc.MustMarshal(&lottery)
	store.Set(GetLotteryIDBytes(lottery.Id), appendedValue)

	// Update Lottery count
	k.SetLotteryCount(ctx, count+1)

	return count
}

// SetLottery sets a specific Lottery in the store
func (k Keeper) SetLottery(ctx sdk.Context, lottery types.Lottery) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LotteryKey))
	b := k.cdc.MustMarshal(&lottery)
	store.Set(GetLotteryIDBytes(lottery.Id), b)
}

// GetLottery returns a Lottery from its id
func (k Keeper) GetLottery(ctx sdk.Context, id uint64) (val types.Lottery, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LotteryKey))

	b := store.Get(GetLotteryIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllLottery returns all lotteries
func (k Keeper) GetAllLottery(ctx sdk.Context) (list []types.Lottery) {
	store := prefix.NewStore(
		ctx.KVStore((k.storeKey)),
		types.KeyPrefix(types.LotteryKey),
	)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Lottery
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetLotteryIDBytes returns the byte representation of the ID
func GetLotteryIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetLotteryIDFromBytes returns ID in uint64 format from a byte array
func GetLotteryIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
