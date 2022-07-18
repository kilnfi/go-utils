package utils

import (
	"fmt"
	"strings"
)

type FlagPrefixer struct {
	prefix, sep string
}

type Option func(fl *FlagPrefixer)

func NewFlagPrefixer(prefix string, opts ...Option) FlagPrefixer {
	fl := FlagPrefixer{
		prefix: SanitizeFlagName(prefix),
		sep:    "-",
	}

	for _, opt := range opts {
		opt(&fl)
	}

	return fl
}

func (fl FlagPrefixer) Prefix() string {
	return fl.prefix
}

func (fl FlagPrefixer) ViperKey(s string) string {
	return ViperKey(fl.FlagName(s))
}

func (fl FlagPrefixer) Env(s string) string {
	return EnvVar(fl.FlagName(s))
}

func (fl FlagPrefixer) FlagName(s string) string {
	return SanitizeFlagName(fmt.Sprintf("%v%v%v", fl.prefix, fl.sep, s))
}

func SeparatorOpt(sep string) Option {
	return func(fl *FlagPrefixer) {
		fl.sep = sep
	}
}

func SanitizeFlagName(flag string) string {
	for _, sep := range []string{
		".", "_", " ",
	} {
		flag = strings.ReplaceAll(flag, sep, "-")
	}

	return strings.ToLower(flag)
}

func EnvVar(flag string) string {
	return strings.ToUpper(strings.ReplaceAll(SanitizeFlagName(flag), "-", "_"))
}

func ViperKey(flag string) string {
	return strings.ReplaceAll(SanitizeFlagName(flag), "-", ".")
}

func FlagDesc(desc, envVar string) string {
	if envVar != "" {
		desc = fmt.Sprintf("%v [env: %v]", desc, envVar)
	}

	return desc
}
