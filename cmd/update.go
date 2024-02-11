package cmd

import (
	"github.com/ASjet/m3/internal/mod"
	"github.com/spf13/cobra"
)

var (
	confirmUpdate bool
)

func init() {
	subCmdUpdate.Flags().BoolVarP(&confirmUpdate,
		"confirm", "y", false, "Confirm download without prompt")
	rootCmd.AddCommand(subCmdUpdate)
}

var subCmdUpdate = &cobra.Command{
	Use:     "update",
	Short:   "Check updates of mods in m3 index",
	Args:    cobra.ExactArgs(0),
	PreRunE: initApiKey,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mod.Update(confirmDownload)
	},
}
