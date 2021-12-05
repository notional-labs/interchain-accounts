package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	ibcatypes "github.com/cosmos/ibc-go/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/interchain-accounts/x/inter-tx/types"
	"github.com/spf13/cobra"
)

const (
	flagPacketTimeoutHeight    = "packet-timeout-height"
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	flagAbsoluteTimeouts       = "absolute-timeouts"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
		getSendTxCmd(),
	)

	return cmd
}

func getSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "send [interchain_account_address] [to_address] [amount] --connection-id",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			cdc := clientCtx.Codec
			ownerAddr := clientCtx.GetFromAddress().String()
			prefixOfDestChain := args[0]
			channelID := args[1]

			listOfFiles := strings.Split(args[2], ",")
			listOfTxType := strings.Split(args[3], ",")

			crossChainMsgs, err := types.SdkMsgsFromFiles(listOfFiles, listOfTxType)
			if err != nil {
				return err
			}
			crossChainMsgsBz, err := ibcatypes.SerializeCosmosTx(
				cdc, crossChainMsgs,
			)
			if err != nil {
				return err
			}

			msg := types.NewMsgSend(
				prefixOfDestChain, ownerAddr, channelID, crossChainMsgsBz,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsConnectionId)

	_ = cmd.MarkFlagRequired(FlagConnectionId)

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// func
