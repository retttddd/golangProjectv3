package cli

import (
	"awesomeProject3/service"
	"awesomeProject3/service/ciphering"
	"awesomeProject3/storage"
	"fmt"
	"github.com/spf13/cobra"
)

var set = &cobra.Command{
	Use:   "set",
	Short: "writes data in",
	Long:  "give 3 parameters: key value password",

	RunE: func(cmd *cobra.Command, args []string) error {
		keys, err := cmd.Flags().GetString("key")
		if err != nil {
			return err
		}
		value, err := cmd.Flags().GetString("value")
		if err != nil {
			return err
		}
		cipherKey, err := cmd.Flags().GetString("cipher-key")
		if err != nil {
			return err
		}
		srv := service.New(storage.NewFsStorage(), ciphering.NewAESEncoder())
		if err := srv.WriteSecret(keys, value, cipherKey); err != nil {
			return err
		}

		fmt.Println("done")
		return nil
	},
}

var get = &cobra.Command{
	Use:   "get",
	Short: "reads data",
	Long:  "give 2 parameters: key password",
	RunE: func(cmd *cobra.Command, args []string) error {
		keys, err := cmd.Flags().GetString("key")
		if err != nil {
			return err
		}
		cipherKey, err := cmd.Flags().GetString("cipher-key")
		if err != nil {
			return err
		}
		srv := service.New(storage.NewFsStorage(), ciphering.NewAESEncoder())
		value, err := srv.ReadSecret(keys, cipherKey)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("decoded data:\n", value)
		return nil
	},
}
