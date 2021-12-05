package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRegisterAccount = "register"
	TypeMsgSend            = "send"
)

var _ sdk.Msg = &MsgSend{}

// NewMsgSend creates a new MsgSend instance
func NewMsgSend(
	prefixOfDestChain string, owner string, channelID string, msgs []byte,
) *MsgSend {
	return &MsgSend{
		PrefixOnDestchain: prefixOfDestChain,
		Owner:             owner,
		ChannelId:         channelID,
		Msgs:              msgs,
	}
}

// Route implements sdk.Msg
func (MsgSend) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgSend) Type() string {
	return TypeMsgSend
}

// GetSigners implements sdk.Msg
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Owner)
	if signer != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// ValidateBasic performs a basic check of the MsgRegisterAccount fields.
func (msg MsgSend) ValidateBasic() error {
	if strings.TrimSpace(msg.PrefixOnDestchain) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing prefix on dest chain")
	}

	if len(msg.Msgs) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty crosschain tx")
	}

	return nil
}

func (msg MsgSend) GetSignBytes() []byte {
	panic("IBC messages do not support amino")
}
