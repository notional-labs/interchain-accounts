package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibcakeeper "github.com/cosmos/ibc-go/modules/apps/27-interchain-accounts/keeper"
	ibcatypes "github.com/cosmos/ibc-go/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	gammtypes "github.com/osmosis-labs/osmosis/x/gamm/types"
	locktypes "github.com/osmosis-labs/osmosis/x/lockup/types"
)

// TrySendCoins builds a banktypes.NewMsgSend and uses the ibc-account module keeper to send the message to another chain
func (keeper Keeper) TrySendCoins(
	ctx sdk.Context,
	owner string,
	prefixOfDestChain string,
	sourceChannelID string,
	msgs []byte,
) error {
	channel, found := keeper.channelkeeper.GetChannel(ctx, ibcatypes.PortID, sourceChannelID)
	if !found {
		return sdkerrors.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", ibcatypes.PortID, sourceChannelID)
	}
	destChannelID := channel.GetCounterparty().GetChannelID()
	destChainAccountAddressBz := ibcakeeper.GenerateAddress(owner + destChannelID + ibcatypes.PortID)
	destChainAccountAddress := sdk.MustBech32ifyAddressBytes(prefixOfDestChain, destChainAccountAddressBz)

	sdkMsgs, err := ibcatypes.DeserializeTx(keeper.cdc, msgs)
	if err != nil {
		return err
	}

	ChangeSignerOfMsgs(sdkMsgs[:], destChainAccountAddress)

	_, err = keeper.iaKeeper.TrySendTx(ctx, owner, sourceChannelID, sdkMsgs)
	return err
}

func ChangeSignerOfMsgs(msgs []sdk.Msg, newSigner string) {
	for id, msg := range msgs {
		switch msg := msg.(type) {
		case *gammtypes.MsgJoinSwapExternAmountIn:
			msg.Sender = newSigner
			msgs[id] = msg

		case *locktypes.MsgLockTokens:
			msg.Owner = newSigner
			msgs[id] = msg
		}
	}
}
