package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// ConfigFile is what it is
	ConfigFile string
)

// InitConfig reads configuration from a file or environment
func InitConfig() {
	// search path
	if ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(ConfigFile)
	} else {
		viper.AddConfigPath("/etc/silphtelescope/")
		viper.AddConfigPath("$HOME/.silphtelescope")
		viper.AddConfigPath(".")
	}

	// env vars
	viper.SetEnvPrefix("SILPHT")
	viper.AutomaticEnv()

	// read config
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
