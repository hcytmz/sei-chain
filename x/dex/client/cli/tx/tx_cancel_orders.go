package tx

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sei-protocol/sei-chain/x/dex/types"
	"github.com/sei-protocol/sei-chain/x/dex/types/utils"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCancelOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-orders [contract address] [cancellations...]",
		Short: "Bulk cancel orders",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argContractAddr := args[0]
			if err != nil {
				return err
			}
			cancellations := []*types.Cancellation{}
			for _, cancellation := range args[1:] {
				newCancel := types.Cancellation{}
				cancelDetails := strings.Split(cancellation, "?")
				argPositionDir, err := utils.GetPositionDirectionFromStr(cancelDetails[0])
				if err != nil {
					return err
				}
				newCancel.PositionDirection = argPositionDir
				argPrice, err := sdk.NewDecFromStr(cancelDetails[1])
				if err != nil {
					return err
				}
				newCancel.Price = argPrice
				newCancel.PriceDenom = cancelDetails[2]
				newCancel.AssetDenom = cancelDetails[3]
				cancellations = append(cancellations, &newCancel)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelOrders(
				clientCtx.GetFromAddress().String(),
				cancellations,
				argContractAddr,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
