package cli

import (
	"awesomeProject3/service"
	"awesomeProject3/service/de_encoding"
	"awesomeProject3/storage"
	"fmt"
	"github.com/spf13/cobra"
)

var set = &cobra.Command{
	Use:   "set",
	Short: "writes data in",
	Long:  "give 3 parameters: key value password",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		srv := service.New(storage.New(), de_encoding.NewAESEncoder(de_encoding.PassToSecretKey(args[2])))
		srv.WriteSecret(args[0], args[1])
		fmt.Println("done")
	},
}

var get = &cobra.Command{
	Use:   "get",
	Short: "reads data",
	Long:  "give 2 parameters: key password",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		srv := service.New(storage.New(), de_encoding.NewAESEncoder(de_encoding.PassToSecretKey(args[1])))
		value, _ := srv.ReadSecret(args[0])
		fmt.Println("decoded data:\n", value)
	},
}
