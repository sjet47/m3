package sub

import (
	"github.com/ASjet/m3/internal/index"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(subCmdInit)
	skipInitApi[subCmdInit.Name()] = true
}

var subCmdInit = &cobra.Command{
	Use:   "init <game_version>",
	Short: "init m3 in current directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return index.Init(args[0])
	},
}
