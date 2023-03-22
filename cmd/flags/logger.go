package flags

import (
	"fmt"

	"github.com/lugondev/signer-key-manager/src/infra/log/zap"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(LevelViperKey, levelDefault)
	_ = viper.BindEnv(LevelViperKey, levelEnv)
	viper.SetDefault(FormatViperKey, formatDefault)
	_ = viper.BindEnv(FormatViperKey, formatEnv)
}

const (
	levelFlag     = "log-level"
	LevelViperKey = "log.level"
	levelDefault  = "info"
	levelEnv      = "LOG_LEVEL"
)

const (
	formatFlag     = "log-format"
	FormatViperKey = "log.format"
	formatDefault  = "text"
	formatEnv      = "LOG_FORMAT"
)

func LoggerFlags(f *pflag.FlagSet) {
	level(f)
	format(f)
}

func level(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log level (one of %q).
Environment variable: %q`, []string{"panic", "error", "warn", "info", "debug"}, levelEnv)
	f.String(levelFlag, levelDefault, desc)
	_ = viper.BindPFlag(LevelViperKey, f.Lookup(levelFlag))
}

func format(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log formatter (one of %q).
Environment variable: %q`, []string{"text", "json"}, formatEnv)
	f.String(formatFlag, formatDefault, desc)
	_ = viper.BindPFlag(FormatViperKey, f.Lookup(formatFlag))
}

func NewLoggerConfig(vipr *viper.Viper) *zap.Config {
	return zap.NewConfig(
		zap.LoggerLevel(vipr.GetString(LevelViperKey)),
		zap.LoggerFormat(vipr.GetString(FormatViperKey)),
	)
}
