package mock

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/ibc/mock/types"
)

type Keeper struct {
	cdc  *codec.Codec
	key  sdk.StoreKey
	port ibc.Port
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, port ibc.Port) Keeper {
	return Keeper{
		cdc:  cdc,
		key:  key,
		port: port,
	}
}

func (k Keeper) GetSequence(ctx sdk.Context, chanid string) (res uint64) {
	store := ctx.KVStore(k.key)
	if store.Has(types.SequenceKey) {
		k.cdc.MustUnmarshalBinaryBare(store.Get(types.SequenceKey), &res)
	} else {
		res = 0
	}

	return
}

func (k Keeper) SetSequence(ctx sdk.Context, chanid string, seq uint64) {
	store := ctx.KVStore(k.key)
	store.Set(types.SequenceKey, k.cdc.MustMarshalBinaryBare(seq))
}

func (k Keeper) CheckAndSetSequence(ctx sdk.Context, chanid string, seq uint64) error {
	if k.GetSequence(ctx, chanid)+1 != seq {
		return errors.New("fjidow;af")
	}
	k.SetSequence(ctx, chanid, seq)
	return nil
}
