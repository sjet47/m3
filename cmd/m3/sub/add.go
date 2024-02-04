package sub

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ASjet/m3/internal/index"
	"github.com/ASjet/m3/internal/mod"
	"github.com/ASjet/m3/internal/util"

	"github.com/spf13/cobra"
)

var (
	optDep    bool
	modLoader string
)

func init() {
	subCmdAdd.Flags().BoolVarP(&optDep,
		"optional", "o", false, "Download optional dependencies")
	subCmdAdd.Flags().StringVarP(&modLoader,
		"modloader", "l", "Forge", "Mod loader")
	rootCmd.AddCommand(subCmdAdd)
}

var subCmdAdd = &cobra.Command{
	Use:   "add <mod_id>...",
	Short: "Add mods to the m3 index",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modIDs, errIndex, err := util.MapErr(strconv.Atoi, args...)
		if err != nil {
			return fmt.Errorf("invalid mod id %q", args[errIndex])
		}

		if err := index.Load(); err != nil {
			return err
		}

		slices.Sort(modIDs)
		slices.Compact(modIDs)

		return mod.Add(modLoader, optDep, modIDs...)
	},
}
