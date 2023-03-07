package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stringer",
	Short: "operator stores and provides ur files secretly",
	Long:  `q`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	rootCmd.AddCommand(get)
	rootCmd.AddCommand(set)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "'%s'", err)
		os.Exit(1)
	}
}
