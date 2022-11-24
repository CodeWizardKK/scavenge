package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	"scavenge/x/scavenge/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSubmitScavenge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-scavenge [solution] [description] [reward]",
		Short: "Broadcast message submit-scavenge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			//正解のハッシュ化
			solution := args[0]
			solutionHash := sha256.Sum256([]byte(solution))
			solutionHashString := hex.EncodeToString(solutionHash[:])
			//description,rewardはそのままセット
			argDescription := args[1]
			argReward := args[2]

			msg := types.NewMsgSubmitScavenge(
				clientCtx.GetFromAddress().String(),
				solutionHashString,
				argDescription,
				argReward,
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
