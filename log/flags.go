package log

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Flags(v *viper.Viper, f *pflag.FlagSet) {
	Level(v, f)
	Format(v, f)
}

const (
	levelFlag     = "log-level"
	LevelViperKey = "log.level"
	levelDefault  = "info"
	levelEnv      = "LOG_LEVEL"
)

func Level(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log level (one of %q).
Environment variable: %q`, []string{"panic", "error", "warn", "info", "debug"}, levelEnv)
	f.String(levelFlag, levelDefault, desc)
	_ = v.BindPFlag(LevelViperKey, f.Lookup(levelFlag))
	v.SetDefault(LevelViperKey, levelDefault)
	_ = v.BindEnv(LevelViperKey, levelEnv)
}

const (
	formatFlag     = "log-format"
	FormatViperKey = "log.format"
	formatDefault  = "text"
	formatEnv      = "LOG_FORMAT"
)

func Format(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log formatter (one of %q).
Environment variable: %q`, []string{"text", "json"}, formatEnv)
	f.String(formatFlag, formatDefault, desc)
	_ = v.BindPFlag(FormatViperKey, f.Lookup(formatFlag))
	v.SetDefault(FormatViperKey, formatDefault)
	_ = v.BindEnv(FormatViperKey, formatEnv)
}

func NewConfigFromViper(v *viper.Viper) *Config {
	return &Config{
		Format: v.GetString(FormatViperKey),
		Level:  v.GetString(LevelViperKey),
	}
}
