package config

import (
	log "github.com/inconshreveable/log15"
	"github.com/spf13/viper"
)

type configuration struct {
	WMAddr      string
	MachineName string
	Debug       bool
}

// Config contains the global configuration
var Config configuration

func init() {
	// set defaults
	viper.RegisterAlias("DEBUG", "debug")
	viper.SetDefault("Debug", false)

	// bind env
	viper.BindEnv("Debug", "DEBUG")
	viper.BindEnv("WMAddr", "WM_ADDR")
	viper.BindEnv("MachineName", "MACHINE_NAME")

	// set config file settings
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	// read configuration file
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("configuration file read")
	}

	// set up logger
	if debug := viper.GetBool("Debug"); debug {
		log.Info("log lvl set to debug")
		log.Root().SetHandler(
			log.LvlFilterHandler(log.LvlDebug, log.StderrHandler),
		)
	} else {
		log.Root().SetHandler(
			log.LvlFilterHandler(log.LvlInfo, log.StderrHandler),
		)
	}

	// set up global configuration object
	Config = configuration{
		WMAddr:      viper.GetString("WMAddr"),
		MachineName: viper.GetString("MachineName"),
		Debug:       viper.GetBool("Debug"),
	}
}
