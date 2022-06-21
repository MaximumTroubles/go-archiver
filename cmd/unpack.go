package cmd

import "github.com/spf13/cobra"

var unpackCmd = &cobra.Command{
	Use:   "unpack",
	Short: "Decompres file command command",
}

func init() {
	rootCmd.AddCommand(unpackCmd)
}
