package cli

import (
	"awesomeProject3/service"
	"awesomeProject3/service/de_encoding"
	"awesomeProject3/storage"
	"fmt"
	"github.com/spf13/cobra"
)

var put = &cobra.Command{
	Use:   "put",
	Short: "writes data in",
	Long:  "give 3 parameters: key value password",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		keys, err := cmd.Flags().GetString("key")
		if err != nil {
			fmt.Println(err)
		}
		value, err := cmd.Flags().GetString("value")
		if err != nil {
			fmt.Println(err)
		}
		pwords, err := cmd.Flags().GetString("cipher-key")
		if err != nil {
			fmt.Println(err)
		}
		p := new(PutData)
		p.Keys = keys
		p.Vals = value
		p.Pword = pwords
		srv := service.New(storage.New(), de_encoding.NewAESEncoder(de_encoding.PassToSecretKey(pwords)))
		srv.WriteSecret(keys, value)
		fmt.Println("done")
	},
}

var get = &cobra.Command{
	Use:   "get",
	Short: "reads data",
	Long:  "give 2 parameters: key password",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		keys, err := cmd.Flags().GetString("key")
		if err != nil {
			fmt.Println(err)
		}
		pwords, err := cmd.Flags().GetString("cipher-key")
		if err != nil {
			fmt.Println(err)
		}
		p := new(GetData)
		p.Keys = keys
		p.Pword = pwords
		srv := service.New(storage.New(), de_encoding.NewAESEncoder(de_encoding.PassToSecretKey(pwords)))
		value, _ := srv.ReadSecret(keys)
		fmt.Println("decoded data:\n", value)
	},
}
