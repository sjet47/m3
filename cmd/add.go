package cmd

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/ASjet/m3/internal/mod"
	"github.com/ASjet/m3/internal/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	confirmDownload bool
	optDep          bool
	modLoader       string
	fileInput       string
)

func init() {
	subCmdAdd.Flags().BoolVarP(&optDep,
		"optional", "o", false, "Download optional dependencies")
	subCmdAdd.Flags().BoolVarP(&confirmDownload,
		"confirm", "y", false, "Confirm download without prompt")
	subCmdAdd.Flags().StringVarP(&modLoader,
		"modloader", "l", "Forge", "Mod loader")
	subCmdAdd.Flags().StringVarP(&fileInput,
		"file", "f", "", "Read mod ids from csv file")
	rootCmd.AddCommand(subCmdAdd)
}

var subCmdAdd = &cobra.Command{
	Use:   "add <mod_id>...",
	Short: "Add mods to the m3 index",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(fileInput) > 0 {
			return cobra.ExactArgs(0)(cmd, args)
		}
		return cobra.MinimumNArgs(1)(cmd, args)
	},
	PreRunE: initApiKey,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(fileInput) > 0 {
			f, err := os.Open(fileInput)
			if err != nil {
				return errors.Wrapf(err, "open file %s error", fileInput)
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				args = append(args, strings.Split(scanner.Text(), ",")[0])
			}
		}

		modIDs, errIndex, err := util.MapErr(strconv.Atoi, args...)
		if err != nil {
			return fmt.Errorf("invalid mod id %q", args[errIndex])
		}

		slices.Sort(modIDs)
		slices.Compact(modIDs)

		return mod.Add(modLoader, confirmDownload, optDep, modIDs...)
	},
}
