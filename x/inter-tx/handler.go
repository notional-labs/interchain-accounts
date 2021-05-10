package inter_tx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	transferTypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/interchainberlin/ica/x/inter-tx/keeper"
	"github.com/interchainberlin/ica/x/inter-tx/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgRegisterAccount:
			res, err := msgServer.Register(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRegisterCommunityAccount:
			res, err := msgServer.RegisterCommunityPool(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSend:
			res, err := msgServer.Send(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleMsgFundProposal(ctx sdk.Context, msg *types.MsgFundProposal, k keeper.Keeper) error {

	senderAddr := k.AuthKeeper.GetModuleAddress("distribution")

	if senderAddr == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message type: %T", msg)
	}

	const timeoutTimestamp = ^uint64(0)

	coin, err := sdk.ParseCoinNormalized(msg.Coin)
	if err != nil {
		return err
	}

	const portId = "transfer"

	//TODO: get from state
	const interchainAccountAddr = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"

	transferMsg := transferTypes.NewMsgTransfer(
		portId, msg.SourceChannel, coin, senderAddr, interchainAccountAddr, clienttypes.ZeroHeight(), timeoutTimestamp,
	)

	_, err = k.TransferKeeper.Transfer(sdk.WrapSDKContext(ctx), transferMsg)

	if err != nil {
		return err
	}

	return nil
}

func handleMsgSendProposal(ctx sdk.Context, msg *types.MsgSendProposal, k keeper.Keeper) error {
	msgServer := keeper.NewMsgServerImpl(k)

	senderAddr := k.AuthKeeper.GetModuleAddress("distribution")

	if senderAddr == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message type: %T", msg)
	}

	sendMsg := types.NewMsgSend(
		"cosmos-sdk",
		msg.SourcePort,
		msg.SourceChannel,
		senderAddr,
		msg.ToAddress,
		msg.Amount,
		msg.Coin,
	)

	const timeoutTimestamp = ^uint64(0)
	_, err := msgServer.Send(sdk.WrapSDKContext(ctx), sendMsg)
	if err != nil {
		return err
	}

	return nil
}

func NewRegisterInterchainAccountProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch msg := content.(type) {
		case *types.MsgSendProposal:
			return handleMsgSendProposal(ctx, msg, k)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
