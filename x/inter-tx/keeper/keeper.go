package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dist "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	ibcacckeeper "github.com/interchainberlin/ica/x/ibc-account/keeper"
	"github.com/interchainberlin/ica/x/inter-tx/types"
)

type Keeper struct {
	cdc      codec.Marshaler
	storeKey sdk.StoreKey
	memKey   sdk.StoreKey

	AuthKeeper types.AccountKeeper
	iaKeeper   ibcacckeeper.Keeper
	distKeeper dist.Keeper
}

func NewKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey, iaKeeper ibcacckeeper.Keeper, distKeeper dist.Keeper, authKeeper types.AccountKeeper) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		distKeeper: distKeeper,
		iaKeeper:   iaKeeper,
		AuthKeeper: authKeeper,
	}
}
