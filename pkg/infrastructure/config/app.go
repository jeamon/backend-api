package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	IsProduction bool   `mapstructure:"is_production"`
	Database     string `mapstructure:"database"`
	Server       struct {
		Port      string `mapstructure:"port"`
		Host      string `mapstructure:"host"`
		CertsFile string `mapstructure:"certs_file"`
		KeyFile   string `mapstructure:"key_file"`
	} `mapstructure:"server"`

	GinDisableReleaseMode bool   `mapstructure:"gin_disable_release_mode"`
	GitCommit             string `mapstructure:"-"`
	GitTag                string `mapstructure:"-"`

	DBPostgresConfig struct {
		Host         string `mapstructure:"host"`
		Port         string `mapstructure:"port"`
		User         string `mapstructure:"user"`
		Password     string `mapstructure:"password"`
		DatabaseName string `mapstructure:"database_name"`
	} `mapstructure:"db_postgres"`

	DBMongoConfig struct {
		Host         string `mapstructure:"host"`
		Port         string `mapstructure:"port"`
		User         string `mapstructure:"user"`
		Password     string `mapstructure:"password"`
		DatabaseName string `mapstructure:"database_name"`
	} `mapstructure:"db_mongo"`
}

func InitConfig(cfgFile string) *Config {
	var config *Config

	if cfgFile != "" {
		// Use config file passed in as argument (From flag)
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home and current directory with name "server.config" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath("app")
		viper.SetConfigName("server.config")
	}

	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database", "postgres")

	var err error
	if err = viper.ReadInConfig(); err != nil {
		panic("Unable to read config: " + err.Error())
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic("unable to decode into struct, " + err.Error())
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err = viper.Unmarshal(&config)
		if err != nil {
			panic("unable to decode into struct: " + err.Error())
		}
	})

	return config
}
