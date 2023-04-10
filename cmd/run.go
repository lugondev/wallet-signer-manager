package cmd

import (
	"fmt"
	"os"

	"github.com/lugondev/wallet-signer-manager/cmd/flags"
	"github.com/lugondev/wallet-signer-manager/pkg/common"
	app "github.com/lugondev/wallet-signer-manager/src"
	"github.com/lugondev/wallet-signer-manager/src/infra/log/zap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRunCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runCmd(cmd, args)
			if err != nil {
				cmd.SilenceUsage = true
			}
			return err
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			preRunBindFlags(viper.GetViper(), cmd.Flags(), "key-manager")
		},
	}

	flags.HTTPFlags(command.Flags())
	flags.ManifestFlags(command.Flags())
	flags.LoggerFlags(command.Flags())
	flags.PGFlags(command.Flags())
	flags.OIDCFlags(command.Flags())
	flags.APIKeyFlags(command.Flags())
	flags.TLSFlags(command.Flags())

	return command
}

func runCmd(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	vipr := viper.GetViper()
	//vipr.AddConfigPath("./")
	vipr.SetConfigFile(".env")
	err := vipr.ReadInConfig()
	if err != nil {
		return err
	}

	cfg, err := flags.NewAppConfig(vipr)
	if err != nil {
		return err
	}
	fmt.Println("app config: ", cfg.ToJson())

	logger, err := zap.NewLogger(cfg.Logger)
	if err != nil {
		return err
	}
	defer syncZapLogger(logger)

	appli, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.WithError(err).Error("could not create app")
		return err
	}

	done := make(chan struct{})
	sig := common.NewSignalListener(func(sig os.Signal) {
		logger.With("sig", sig.String()).Warn("signal intercepted")
		if err = appli.Stop(ctx); err != nil {
			logger.WithError(err).Error("application stopped with errors")
		}
		close(done)
	})

	defer sig.Close()

	err = appli.Start(ctx)
	if err != nil {
		logger.WithError(err).Error("application failed to start")
		return err
	}

	<-done

	return nil
}
