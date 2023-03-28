package cli

import (
	"awesomeProject3/rest"
	"awesomeProject3/service"
	"awesomeProject3/service/ciphering"
	"awesomeProject3/storage"
	"crypto/rand"
	"fmt"
	"github.com/spf13/cobra"
)

type constReader struct {
}

func (r *constReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}

	return len(p), nil
}

var set = &cobra.Command{
	Use:   "set",
	Short: "writes data in",
	Long:  "give 3 parameters: key value password",

	RunE: func(cmd *cobra.Command, args []string) error {
		cReader := &constReader{}
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
		srv := service.New(storage.NewFsStorage("./data/test.json"),
			ciphering.NewAESEncoder(ciphering.NewRandomNonceProducer(rand.Reader)),
			ciphering.NewAESEncoder(ciphering.NewRandomNonceProducer(cReader)))
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
		cReader := &constReader{}
		keys, err := cmd.Flags().GetString("key")
		if err != nil {
			return err
		}
		cipherKey, err := cmd.Flags().GetString("cipher-key")
		if err != nil {
			return err
		}
		srv := service.New(storage.NewFsStorage("./data/test.json"),
			ciphering.NewAESEncoder(ciphering.NewRandomNonceProducer(rand.Reader)),
			ciphering.NewAESEncoder(ciphering.NewRandomNonceProducer(cReader)))
		value, err := srv.ReadSecret(keys, cipherKey)
		if err != nil {
			return err
		}
		fmt.Println("decoded data:\n", value)
		return nil
	},
}

var server = &cobra.Command{
	Use:   "server",
	Short: "starts server",
	Long:  "give 2 parameters: port filepath",
	Run: func(cmd *cobra.Command, args []string) {
		cReader := &constReader{}
		secretService := service.New(storage.NewFsStorage("./data/test.json"),
			ciphering.NewAESEncoder(ciphering.NewRandomNonceProducer(rand.Reader)),
			ciphering.NewAESEncoder(ciphering.NewRandomNonceProducer(cReader)))
		srv := rest.NewHttpServer(secretService)
		srv.Start()
	},
}
