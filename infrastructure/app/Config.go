package app

import (
	"log"

	"github.com/spf13/viper"
)

type ApplicationConfig struct {
	DatabaseConnection string
}

// ConfigSetup will prepare and setup the viper instance to the correct config file.
func ConfigSetup(configName, configPath string) {
	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath(configPath)
}

// GetConfig initializes the configuration instance to the values described in the config.toml file.
func GetConfig() *ApplicationConfig {
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("fatal error config file: %s", err)
	}

	validateVariablesAreSet([]string{
		"DatabaseConnection",
	})

	return &ApplicationConfig{}
}

// validateVariablesAreSet will assert the existence of each variable,
// and kill the application when a wanted variable does not exist in the config.
func validateVariablesAreSet(variables []string) {
	for i := range variables {
		if !viper.IsSet(variables[i]) {
			log.Fatalf("%s variable was not set!\nAborting application start!", variables[i])
		}
	}
}
