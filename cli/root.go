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
	rootCmd.AddCommand(server)
	rootCmd.AddCommand(healthcheck)
	healthcheck.Flags().StringP("port", "o", "", "port you use to start server")
	server.Flags().StringP("port", "o", "", "port you use to start server")
	server.Flags().StringP("path", "a", "", "path you use to store data")
	server.Flags().StringP("database", "x", "", "add url to your database")
	set.Flags().StringP("key", "k", "", "key you use to identify data")
	set.Flags().StringP("value", "v", "", "data you pass ")
	set.Flags().StringP("cipher-key", "p", "", "cipher-key you use to make your data secret ")
	set.Flags().StringP("path", "a", "", "path you use to store data")
	get.Flags().StringP("key", "k", "", "key you use to identify data")
	get.Flags().StringP("cipher-key", "p", "", "cipher-key you use to make your data secret ")
	get.Flags().StringP("path", "a", "", "path you use to store data")
	return rootCmd.Execute()
}
