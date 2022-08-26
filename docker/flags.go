package docker

import (
	kilncmd "github.com/kilnfi/go-utils/cmd/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	hostFlag     = "docker-host"
	hostViperKey = "docker.host"
	hostEnv      = "DOCKER_HOST"
)

// HostFlag register flag for docker host
func HostFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := kilncmd.FlagDesc("Docker host", hostEnv)
	f.String(hostFlag, "", desc)
	_ = v.BindPFlag(hostViperKey, f.Lookup(hostFlag))
	_ = v.BindEnv(hostViperKey, hostEnv)
}

func GetHost(v *viper.Viper) string {
	return v.GetString(hostViperKey)
}

const (
	apiVersionFlag     = "docker-api-version"
	apiVersionViperKey = "docker.api-version"
	apiVersionEnv      = "DOCKER_API_VERSION"
)

// APIVersionFlag register flag for docker API version
func APIVersionFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := kilncmd.FlagDesc("Docker API version", apiVersionEnv)
	f.String(apiVersionFlag, "", desc)
	_ = v.BindPFlag(apiVersionViperKey, f.Lookup(apiVersionFlag))
	_ = v.BindEnv(apiVersionViperKey, apiVersionEnv)
}

func GetAPIVersion(v *viper.Viper) string {
	return v.GetString(apiVersionViperKey)
}

func ClientConfigFromViper(v *viper.Viper) *ClientConfig {
	return &ClientConfig{
		Host:    GetHost(v),
		Version: GetAPIVersion(v),
	}
}
