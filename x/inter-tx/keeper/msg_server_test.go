package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"

	"github.com/cosmos/interchain-accounts/v7/x/inter-tx/keeper"
	"github.com/cosmos/interchain-accounts/v7/x/inter-tx/types"
)

func (s *KeeperTestSuite) TestRegisterInterchainAccount() {
	var (
		owner string
		path  *ibctesting.Path
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success", func() {}, true,
		},
		{
			"failure - port is already bound",
			func() {
				GetICAApp(s.chainA).IBCKeeper.PortKeeper.BindPort(s.chainA.GetContext(), TestPortID)
			},
			false,
		},
		{
			"faliure - owner is empty",
			func() {
				owner = ""
			},
			false,
		},
		{
			"failure - channel is already active",
			func() {
				portID, err := icatypes.NewControllerPortID(owner)
				s.Require().NoError(err)

				channel := channeltypes.NewChannel(
					channeltypes.OPEN,
					channeltypes.ORDERED,
					channeltypes.NewCounterparty(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID),
					[]string{path.EndpointA.ConnectionID},
					path.EndpointA.ChannelConfig.Version,
				)

				GetICAApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), portID, ibctesting.FirstChannelID, channel)
				GetICAApp(s.chainA).ICAControllerKeeper.SetActiveChannelID(s.chainA.GetContext(), ibctesting.FirstConnectionID, portID, ibctesting.FirstChannelID)
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			s.SetupTest()

			owner = TestOwnerAddress // must be explicitly changed

			path = NewICAPath(s.chainA, s.chainB)
			s.coordinator.SetupConnections(path)

			tc.malleate() // malleate mutates test data

			msgSrv := keeper.NewMsgServerImpl(GetICAApp(s.chainA).InterTxKeeper)
			msg := types.NewMsgRegisterAccount(owner, path.EndpointA.ConnectionID, path.EndpointA.ChannelConfig.Version)

			res, err := msgSrv.RegisterAccount(sdk.WrapSDKContext(s.chainA.GetContext()), msg)

			if tc.expPass {
				s.Require().NoError(err)
				s.Require().NotNil(res)
			} else {
				s.Require().Error(err)
				s.Require().Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestSubmitTx() {
	var (
		path                      *ibctesting.Path
		registerInterchainAccount bool
		owner                     string
		connectionID              string
		icaMsg                    sdk.Msg
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionID = path.EndpointA.ConnectionID
			}, true,
		},
		{
			"failure - owner address is empty", func() {
				registerInterchainAccount = true
				owner = ""
				connectionID = path.EndpointA.ConnectionID
			}, false,
		},
		{
			"failure - active channel does not exist for connection ID", func() {
				registerInterchainAccount = true
				owner = TestOwnerAddress
				connectionID = "connection-100"
			}, false,
		},
		{
			"failure - active channel does not exist for port ID", func() {
				registerInterchainAccount = true
				owner = "cosmos153lf4zntqt33a4v0sm5cytrxyqn78q7kz8j8x5"
				connectionID = path.EndpointA.ConnectionID
			}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			s.SetupTest()

			icaAppA := GetICAApp(s.chainA)
			icaAppB := GetICAApp(s.chainB)

			path = NewICAPath(s.chainA, s.chainB)
			s.coordinator.SetupConnections(path)

			tc.malleate() // malleate mutates test data

			if registerInterchainAccount {
				err := SetupICAPath(path, TestOwnerAddress)
				s.Require().NoError(err)

				portID, err := icatypes.NewControllerPortID(TestOwnerAddress)
				s.Require().NoError(err)

				// Get the address of the interchain account stored in state during handshake step
				interchainAccountAddr, found := GetICAApp(s.chainA).ICAControllerKeeper.GetInterchainAccountAddress(s.chainA.GetContext(), path.EndpointA.ConnectionID, portID)
				s.Require().True(found)

				icaAddr, err := sdk.AccAddressFromBech32(interchainAccountAddr)
				s.Require().NoError(err)

				// Check if account is created
				interchainAccount := icaAppB.AccountKeeper.GetAccount(s.chainB.GetContext(), icaAddr)
				s.Require().Equal(interchainAccount.GetAddress().String(), interchainAccountAddr)

				// Create bank transfer message to execute on the host
				icaMsg = &banktypes.MsgSend{
					FromAddress: interchainAccountAddr,
					ToAddress:   s.chainB.SenderAccount.GetAddress().String(),
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				}
			}

			msgSrv := keeper.NewMsgServerImpl(icaAppA.InterTxKeeper)
			msg, err := types.NewMsgSubmitTx(icaMsg, connectionID, owner)
			s.Require().NoError(err)

			res, err := msgSrv.SubmitTx(sdk.WrapSDKContext(s.chainA.GetContext()), msg)

			if tc.expPass {
				s.Require().NoError(err)
				s.Require().NotNil(res)
			} else {
				s.Require().Error(err)
				s.Require().Nil(res)
			}
		})
	}
}
