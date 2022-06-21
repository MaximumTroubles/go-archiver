package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "Go Archiver",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		handleErr(err)
	}
}

// Not the best pratice for handle error but we simplify to learn new material
func handleErr(err error) {
	_, _ = fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}
