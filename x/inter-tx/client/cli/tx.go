package cli

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/interchainberlin/ica/x/inter-tx/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		getRegisterAccountCmd(),
		getSendTxCmd(),
		GetCmdSubmitProposal(),
	)

	return cmd
}

func getRegisterAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "register",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sourcePort := viper.GetString(FlagSourcePort)
			sourceChannel := viper.GetString(FlagSourceChannel)

			msg := types.NewMsgRegisterAccount(
				sourcePort,
				sourceChannel,
				clientCtx.GetFromAddress().String(),
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(fsSourcePort)
	cmd.Flags().AddFlagSet(fsSourceChannel)

	_ = cmd.MarkFlagRequired(FlagSourcePort)
	_ = cmd.MarkFlagRequired(FlagSourceChannel)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "send [type] [to_address] [amount]",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			chainType := args[0]
			fromAddress := clientCtx.GetFromAddress()
			toAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			sourcePort := viper.GetString(FlagSourcePort)
			sourceChannel := viper.GetString(FlagSourceChannel)

			msg := types.NewMsgSend(
				chainType,
				sourcePort,
				sourceChannel,
				fromAddress,
				toAddress,
				amount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(fsSourcePort)
	cmd.Flags().AddFlagSet(fsSourceChannel)

	_ = cmd.MarkFlagRequired(FlagSourcePort)
	_ = cmd.MarkFlagRequired(FlagSourceChannel)

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

const HELPER_TEXT = `Submit a community pool register proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal community-pool-register <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "I",
  "description": "Registering Interchain Account",
  "sourcePort": "ibcaccount",
  "sourceChannel": "1000stake",
  "deposit": "1000stake",
}
`

func GetCmdSubmitProposal() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "community-pool-register [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a community pool register interchain-account proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(HELPER_TEXT, version.AppName, bech32PrefixAccAddr),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := ParseCommunityPoolSpendProposalWithDeposit(clientCtx.JSONMarshaler, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			deposit, err := sdk.ParseCoinsNormalized("50stake")

			content := types.NewMsgRegisterProposal(proposal.Title, proposal.Description, proposal.SourcePort, proposal.SourceChannel)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func ParseCommunityPoolSpendProposalWithDeposit(cdc codec.JSONMarshaler, proposalFile string) (types.MsgRegisterProposal, error) {
	proposal := types.MsgRegisterProposal{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
