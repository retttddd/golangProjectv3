package cli

import (
	"github.com/spf13/cobra"
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
	set.Flags().StringP("key", "k", "", "key you use to identify data")
	set.Flags().StringP("value", "v", "", "data you pass ")
	set.Flags().StringP("cipher-key", "p", "", "cipher-key you use to make your data secret ")
	get.Flags().StringP("key", "k", "", "key you use to identify data")
	get.Flags().StringP("cipher-key", "p", "", "cipher-key you use to make your data secret ")
	return rootCmd.Execute()
}
