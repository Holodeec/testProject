package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server   Server
	Database DatabaseConf
	Cache    Cache
}

func Get() *Configuration {
	conf := &Configuration{}
	viper.AddConfigPath("./configs")
	viper.SetConfigName("local_config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()
	viper.BindEnv("cache.ttl", "CACHE_TTL")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}

	return conf
}

type DatabaseConf struct {
	Dialect  string `yaml:"dialect"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
}

type Server struct {
	Port    string
	Mode    string
	Version string
}

type Cache struct {
	Ttl time.Duration `ENV:"CACHE_TTL"`
}

func (s Server) GetAddr() string {
	return fmt.Sprintf(":%s", s.Port)
}

func (d *DatabaseConf) GetConnectionString() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		d.Dialect, d.Username, d.Password, d.Host, d.Port, d.Database)

}
