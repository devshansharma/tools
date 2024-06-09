package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// New for taking a object and a filepath and fill the object with the provided json configuration
func New(cfg interface{}, jsonFilePath string) interface{} {
	viper.SetConfigFile(jsonFilePath)
	// viper.SetConfigName(filename)
	// viper.SetConfigType("json")
	// viper.AddConfigPath("./configs")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error read config file: %w", err))
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal to app config: %w", err))
	}

	return cfg
}
