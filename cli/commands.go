package cli

import (
	health "awesomeProject3/healthcheck"
	"awesomeProject3/rest"
	"awesomeProject3/service"
	"awesomeProject3/service/ciphering"
	"awesomeProject3/storage"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"syscall"
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
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}
		srv := service.New(storage.NewFsStorage(path),
			ciphering.NewAESEncoder(ciphering.NewRandomNonceProducer(rand.Reader)),
			ciphering.NewAESEncoder(ciphering.NewRandomNonceProducer(cReader)))
		if err := srv.WriteSecret(keys, &service.SecretServiceModel{Value: &value}, cipherKey); err != nil {
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
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}

		srv := service.New(storage.NewFsStorage(path),
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
var healthcheck = &cobra.Command{
	Use:   "healthcheck",
	Short: "checks health",
	Long:  "give 2 parameters: key password",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			return err
		}
		return health.Check(port)
	},
}

var server = &cobra.Command{
	Use:   "server",
	Short: "starts server",
	Long:  "give 2 parameters: port and filepath",
	RunE: func(cmd *cobra.Command, args []string) error {
		cReader := &constReader{}
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			return err
		}
		var secretService service.SecretService
		dbUrl, err := cmd.Flags().GetString("database")
		if err != nil {
			return err
		}
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}
		if dbUrl != "" && path == "" {
			secretService = service.New(storage.NewDbStorage(dbUrl),
				ciphering.NewAESEncoder(ciphering.NewRandomNonceProducer(rand.Reader)),
				ciphering.NewAESEncoder(ciphering.NewRandomNonceProducer(cReader)))
		} else if dbUrl == "" && path != "" {
			secretService = service.New(storage.NewFsStorage(path),
				ciphering.NewAESEncoder(ciphering.NewRandomNonceProducer(rand.Reader)),
				ciphering.NewAESEncoder(ciphering.NewRandomNonceProducer(cReader)))
		} else {
			return errors.New("inconsistend parametr list")
		}
		srv := rest.NewSecretRestAPI(secretService, port)
		serverCtx, serverCancel := context.WithCancel(cmd.Context())
		go func() {
			sign := make(chan os.Signal, 1)
			signal.Notify(sign, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			defer signal.Stop(sign)
			defer serverCancel()

			select {
			case <-sign:
			case <-serverCtx.Done():
			}
		}()

		err = srv.Start(serverCtx)

		log.Println("Done")
		return err
	},
}
