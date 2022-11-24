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

func CmdCommitSolution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-solution [solution]",
		Short: "Broadcast message commit-solution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			//正解のハッシュ化
			solution := args[0]
			solutionHash := sha256.Sum256([]byte(solution))
			solutionHashString := hex.EncodeToString(solutionHash[:])

			//正解 + 正解を設定したアカウントをハッシュ化
			var scavenger = clientCtx.GetFromAddress().String()
			var solutionScavengerHash = sha256.Sum256([]byte(solution + scavenger))
			var solutionScavengerHashString = hex.EncodeToString(solutionScavengerHash[:])

			msg := types.NewMsgCommitSolution(
				clientCtx.GetFromAddress().String(),
				solutionHashString,
				solutionScavengerHashString,
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
