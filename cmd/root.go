/*
Copyright © 2024 Aryan Sjet <sjet@asjet.dev>
*/
package cmd

import (
	"errors"
	"os"

	"github.com/ASjet/m3/internal/index"
	"github.com/ASjet/m3/internal/mod"

	"github.com/spf13/cobra"
)

var (
	cfApiKey string
	skipLoad = make(map[string]bool)
	skipSave = make(map[string]bool)
)

func initApiKey(cmd *cobra.Command, args []string) error {
	if len(cfApiKey) == 0 {
		if cfApiKey = os.Getenv("CURSE_FORGE_APIKEY"); len(cfApiKey) == 0 {
			return errors.New("no CurseForge API key provided")
		}
	}
	mod.Init(cfApiKey)
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "m3",
	Short: "A Minecraft Mod Manager (https://github.com/ASjet/m3)",
	Args:  cobra.ExactArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if skipLoad[cmd.Name()] {
			return nil
		}
		return index.Load()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if skipSave[cmd.Name()] {
			return nil
		}
		return index.Save()
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
