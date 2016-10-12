package main

import (
	"fmt"
	"log"

	"github.com/mievstac/QuickJob/config"
	"github.com/mievstac/QuickJob/database"
	"github.com/mievstac/QuickJob/validators"
)

type SetupConfig struct {
	Config *config.Config
}

func preSetup() *SetupConfig {
	setupConf := SetupConfig{}
	conf, err := config.OpenConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	if conf != nil {
		setupConf.Config = conf
	}
	return &setupConf
}

func main() {
	var k int
	k++
	k++
	fmt.Println(k)
	fmt.Println("### QuickJobs")
	setup := preSetup()
	redisConnection := database.NewRedisConn(setup.Config)
	check, err := validators.Validate(redisConnection, setup.Config)
	if err != nil {
		log.Fatal(err)
	}
	if check == true {
		fmt.Println("ALL GOOD")
	}
}
