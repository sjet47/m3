package sub

import (
	"fmt"

	"github.com/ASjet/m3/internal/index"
	"github.com/spf13/cobra"
)

func init() {
	skipSave[subCmdInit.Name()] = true
	rootCmd.AddCommand(subCmdLs)
}

var subCmdLs = &cobra.Command{
	Use:   "ls",
	Short: "List mods in m3 index",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(index.Mods.String())
	},
}
