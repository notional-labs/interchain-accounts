package inter_tx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
		case *types.MsgSend:
			res, err := msgServer.Send(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleMsgRegisterProposal(ctx sdk.Context, msg *types.MsgRegisterProposal, k keeper.Keeper) error {
	senderAddr := k.AuthKeeper.GetModuleAddress("distribution")

	if senderAddr == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message type: %T", msg)
	}
	// Define custom logic here for registering an account on behalf of the blockcahin
	//msgRegister := types.NewMsgRegisterAccount(msg.SourcePort, msg.SourceChannel, senderAddr.String())

	//	msgServer := keeper.NewMsgServerImpl(k)
	//	_, err := msgServer.Register(sdk.WrapSDKContext(ctx), msgRegister)
	err := k.RegisterIBCAccount(
		ctx,
		senderAddr,
		msg.SourcePort,
		msg.SourceChannel,
	)
	if err != nil {
		return err
	}

	return nil
}

func NewRegisterInterchainAccountProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch msg := content.(type) {
		case *types.MsgRegisterProposal:
			return handleMsgRegisterProposal(ctx, msg, k)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
