package cli

import (
	"github.com/spf13/cobra"
)

var (
	Keys   string
	Vals   string
	Pwords string
)

var rootCmd = &cobra.Command{
	Use:   "secret",
	Short: "operator stores and provides ur files secretly",
	Long:  `q`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() error {
	rootCmd.AddCommand(get)
	rootCmd.AddCommand(set)
	set.Flags().StringVarP(&Keys, "key", "k", "", "key you use to identify data")
	set.Flags().StringVarP(&Vals, "value", "v", "", "data you pass ")
	set.Flags().StringVarP(&Pwords, "cipher-key", "p", "", "cipher-key you use to make your data secret ")
	get.Flags().StringVarP(&Keys, "key", "k", "", "key you use to identify data")
	get.Flags().StringVarP(&Pwords, "cipher-key", "p", "", "cipher-key you use to make your data secret ")
	return rootCmd.Execute()
}
