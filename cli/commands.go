package cli

import (
	"awesomeProject3/service"
	"awesomeProject3/service/ciphering"
	"awesomeProject3/storage"
	"fmt"
	"github.com/spf13/cobra"
)

type setData struct {
	Keys  string
	Vals  string
	Pword string
}

type getData struct {
	Keys  string
	Pword string
}

var set = &cobra.Command{
	Use:   "set",
	Short: "writes data in",
	Long:  "give 3 parameters: key value password",

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
		p := new(setData)
		p.Keys = keys
		p.Vals = value
		p.Pword = pwords
		srv := service.New(storage.NewFsStorage(), ciphering.NewAESEncoder())
		err1 := srv.WriteSecret(keys, value, pwords)
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		fmt.Println("done")
	},
}

var get = &cobra.Command{
	Use:   "get",
	Short: "reads data",
	Long:  "give 2 parameters: key password",
	Run: func(cmd *cobra.Command, args []string) {
		keys, err := cmd.Flags().GetString("key")
		if err != nil {
			fmt.Println(err)
		}
		pwords, err := cmd.Flags().GetString("cipher-key")
		if err != nil {
			fmt.Println(err)
		}
		p := new(getData)
		p.Keys = keys
		p.Pword = pwords
		srv := service.New(storage.NewFsStorage(), ciphering.NewAESEncoder())
		value, err := srv.ReadSecret(keys, pwords)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("decoded data:\n", value)
	},
}
