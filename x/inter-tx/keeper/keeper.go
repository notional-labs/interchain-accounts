package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	ibcacckeeper "github.com/cosmos/ibc-go/modules/apps/27-interchain-accounts/keeper"
	channelkeeper "github.com/cosmos/ibc-go/modules/core/04-channel/keeper"
)

type Keeper struct {
	cdc           codec.Codec
	iaKeeper      ibcacckeeper.Keeper
	channelkeeper channelkeeper.Keeper
}

func NewKeeper(cdc codec.Codec, iaKeeper ibcacckeeper.Keeper, channelkeeper channelkeeper.Keeper) Keeper {
	return Keeper{
		cdc:           cdc,
		channelkeeper: channelkeeper,
		iaKeeper:      iaKeeper,
	}
}
