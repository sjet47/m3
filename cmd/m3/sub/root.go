/*
Copyright © 2024 Aryan Sjet <sjet@asjet.dev>
*/
package sub

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "m3",
	Short: "A Minecraft Mod Manager (https://github.com/ASjet/m3)",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
