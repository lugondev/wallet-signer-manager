package cmd

import (
	"strings"

	"github.com/lugondev/signer-key-manager/src/infra/log/zap"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// NewCommand create root command
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "key-manager",
		Short: "Run Quorum Key Manager",
	}

	rootCmd.AddCommand(newRunCommand())
	rootCmd.AddCommand(newMigrateCommand())

	return rootCmd
}

func syncZapLogger(logger *zap.Logger) {
	_ = logger.Sync()
}

// TODO: Use Viper to parse the flags instead of manually doing it
func preRunBindFlags(vipr *viper.Viper, flgs *pflag.FlagSet, ignore string) {
	for _, vk := range vipr.AllKeys() {
		if ignore != "" && strings.HasPrefix(vk, ignore) {
			continue
		}

		// Convert viperKey to cmd flag name
		// For example: 'rest.api' to "rest-api"
		name := strings.Replace(vk, ".", "-", -1)

		// Only bind in case command flags contain the name
		if flgs.Lookup(name) != nil {
			_ = viper.BindPFlag(vk, flgs.Lookup(name))
		}
	}
}
