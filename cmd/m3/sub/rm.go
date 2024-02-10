package sub

import (
	"fmt"
	"strconv"

	"github.com/ASjet/m3/internal/index"
	"github.com/ASjet/m3/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(subCmdRm)
}

var subCmdRm = &cobra.Command{
	Use:   "rm <mod_id>...",
	Short: "Remove mods from m3 index",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modIDs, errIndex, err := util.MapErr(strconv.Atoi, args...)
		if err != nil {
			return fmt.Errorf("invalid mod id %q", args[errIndex])
		}
		return index.Remove(modIDs...)
	},
}
