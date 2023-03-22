package cmd

import (
	"fmt"
	"net/url"

	"github.com/lugondev/signer-key-manager/src/infra/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // nolint
	_ "github.com/golang-migrate/migrate/v4/source/file"       // nolint
	"github.com/lugondev/signer-key-manager/cmd/flags"
	"github.com/lugondev/signer-key-manager/src/infra/log/zap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMigrateCommand() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run migration",
		PreRun: func(cmd *cobra.Command, args []string) {
			preRunBindFlags(viper.GetViper(), cmd.Flags(), "key-manager")
		},
	}

	// Register Up command
	upCmd := &cobra.Command{
		Use:   "up [target]",
		Short: "Executes all migrations",
		Long:  "Executes all available migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := migrateUp(); err != nil {
				cmd.SilenceUsage = true
				return err
			}
			return nil
		},
	}
	migrateCmd.AddCommand(upCmd)
	flags.LoggerFlags(upCmd.Flags())
	flags.PGFlags(upCmd.Flags())

	// Register Down command
	downCmd := &cobra.Command{
		Use:   "down",
		Short: "Reverts last migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := migrateDown(); err != nil {
				cmd.SilenceUsage = true
				return err
			}
			return nil
		},
	}
	migrateCmd.AddCommand(downCmd)
	flags.LoggerFlags(downCmd.Flags())
	flags.PGFlags(downCmd.Flags())

	// Register Reset command
	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reverts all migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := migrateReset(); err != nil {
				cmd.SilenceUsage = true
				return err
			}
			return nil
		},
	}
	migrateCmd.AddCommand(resetCmd)
	flags.LoggerFlags(resetCmd.Flags())
	flags.PGFlags(resetCmd.Flags())

	return migrateCmd
}

func migrateUp() error {
	vipr := viper.GetViper()
	logger, err := initLogger(vipr)
	if err != nil {
		return err
	}

	m, err := initMigrations(vipr, logger)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		errMessage := "failed to execute migrations"
		logger.WithError(err).Error(errMessage)
		return err
	}

	logger.Info("migrations executed successfully")
	return nil
}

func migrateDown() error {
	vipr := viper.GetViper()
	logger, err := initLogger(vipr)
	if err != nil {
		return err
	}

	m, err := initMigrations(vipr, logger)
	if err != nil {
		return err
	}

	err = m.Steps(-1)
	if err != nil && err != migrate.ErrNoChange {
		errMessage := "failed to downgrade migrations"
		logger.WithError(err).Error(errMessage)
		return err
	}

	logger.Info("migration down successfully")
	return nil
}

func migrateReset() error {
	vipr := viper.GetViper()
	logger, err := initLogger(vipr)
	if err != nil {
		return err
	}

	m, err := initMigrations(vipr, logger)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		errMessage := "failed to reset all migrations"
		logger.WithError(err).Error(errMessage)
		return err
	}

	logger.Info("migrations reset successfully")
	return nil
}

func initMigrations(vipr *viper.Viper, logger log.Logger) (*migrate.Migrate, error) {
	pgCfg := flags.NewPostgresConfig(vipr)
	var userPass string
	if pgCfg.Password != "" {
		userPass = fmt.Sprintf("%s:%s", url.QueryEscape(pgCfg.User), url.QueryEscape(pgCfg.Password))
	} else {
		userPass = url.QueryEscape(pgCfg.User)
	}

	dbURL := fmt.Sprintf("postgres://%s@%s:%s/%s?", userPass, pgCfg.Host, pgCfg.Port, pgCfg.Database)
	params := url.Values{}
	params.Add("sslmode", pgCfg.SSLMode)
	params.Add("sslcert", pgCfg.TLSCert)
	params.Add("sslkey", pgCfg.TLSKey)
	params.Add("sslrootcert", pgCfg.TLSCA)

	m, err := migrate.New("file:///migrations", dbURL+params.Encode())
	if err != nil {
		errMessage := "failed to create migration instance"
		logger.WithError(err).Error(errMessage)
		return nil, err
	}

	return m, nil
}

func initLogger(vipr *viper.Viper) (log.Logger, error) {
	logCfg := flags.NewLoggerConfig(vipr)
	return zap.NewLogger(logCfg)
}
