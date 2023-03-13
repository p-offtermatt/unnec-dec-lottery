package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"lottery/x/lottery/types"
)

// GetLotteryPotsCount get the total number of lotteryPots
func (k Keeper) GetLotteryPotsCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.LotteryPotsCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetLotteryPotsCount set the total number of lotteryPots
func (k Keeper) SetLotteryPotsCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.LotteryPotsCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendLotteryPots appends a lotteryPots in the store with a new id and update the count
func (k Keeper) AppendLotteryPots(
	ctx sdk.Context,
	lotteryPots types.LotteryPots,
) uint64 {
	// Create the lotteryPots
	count := k.GetLotteryPotsCount(ctx)

	// Set the ID of the appended value
	lotteryPots.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LotteryPotsKey))
	appendedValue := k.cdc.MustMarshal(&lotteryPots)
	store.Set(GetLotteryPotsIDBytes(lotteryPots.Id), appendedValue)

	// Update lotteryPots count
	k.SetLotteryPotsCount(ctx, count+1)

	return count
}

// SetLotteryPots set a specific lotteryPots in the store
func (k Keeper) SetLotteryPots(ctx sdk.Context, lotteryPots types.LotteryPots) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LotteryPotsKey))
	b := k.cdc.MustMarshal(&lotteryPots)
	store.Set(GetLotteryPotsIDBytes(lotteryPots.Id), b)
}

// GetLotteryPots returns a lotteryPots from its id
func (k Keeper) GetLotteryPots(ctx sdk.Context, id uint64) (val types.LotteryPots, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LotteryPotsKey))
	b := store.Get(GetLotteryPotsIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveLotteryPots removes a lotteryPots from the store
func (k Keeper) RemoveLotteryPots(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LotteryPotsKey))
	store.Delete(GetLotteryPotsIDBytes(id))
}

// GetAllLotteryPots returns all lotteryPots
func (k Keeper) GetAllLotteryPots(ctx sdk.Context) (list []types.LotteryPots) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LotteryPotsKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LotteryPots
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetLotteryPotsIDBytes returns the byte representation of the ID
func GetLotteryPotsIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetLotteryPotsIDFromBytes returns ID in uint64 format from a byte array
func GetLotteryPotsIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
