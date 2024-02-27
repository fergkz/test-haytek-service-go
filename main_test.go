package main

import (
	"log"
	"testing"
)

func Test(t *testing.T) {
	config := new(Config)
	config.Load("config-test.yml")

	log.Println(config)
}
