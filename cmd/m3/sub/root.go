/*
Copyright © 2024 Aryan Sjet <sjet@asjet.dev>
*/
package sub

import (
	"errors"
	"m3/internal/mod"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfApiKey    string
	skipInitApi = make(map[string]bool)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "m3",
	Short: "A Minecraft Mod Manager (https://github.com/ASjet/m3)",
	Args:  cobra.ExactArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if skipInitApi[cmd.Name()] {
			return nil
		}

		if len(cfApiKey) == 0 {
			cfApiKey = os.Getenv("CURSE_FORGE_APIKEY")
		}

		if len(cfApiKey) == 0 {
			return errors.New("no CurseForge API key provided")
		}

		mod.Init(cfApiKey)

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfApiKey, "api-key", "k", "",
		"CurseForge API Key, or use env CURSE_FORGE_APIKEY if not set")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
