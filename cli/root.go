package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type PutData struct {
	Keys  string
	Vals  string
	Pword string
}

type GetData struct {
	Keys  string
	Pword string
}

var (
	Keys   string
	Vals   string
	Pwords string
)

var rootCmd = &cobra.Command{
	Use:   "secret",
	Short: "operator stores and provides ur files secretly",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	rootCmd.AddCommand(put)
	rootCmd.AddCommand(get)
	put.Flags().StringVar(&Keys, "key", "", "key you use to identify data")
	put.Flags().StringVar(&Vals, "value", "", "data you pass ")
	put.Flags().StringVar(&Pwords, "cipher-key", "", "cipher-key you use to make your data secret ")
	get.Flags().StringVar(&Keys, "key", "", "key you use to identify data")
	get.Flags().StringVar(&Pwords, "cipher-key", "", "cipher-key you use to make your data secret ")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "'%s'", err)
		os.Exit(1)
	}
}
