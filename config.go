package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
}

func (config *Config) Load(Filename string) error {

	splitData := strings.Split(Filename, ".")

	viper.SetConfigName(splitData[0])
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType(splitData[1])

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		return err
	}

	viper.Unmarshal(config)

	return nil
}
