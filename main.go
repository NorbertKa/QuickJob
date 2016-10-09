package main

import (
	"fmt"
	"log"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/mievstac/QuickJob/config"
	"github.com/mievstac/QuickJob/database"
	"github.com/mievstac/QuickJob/validator"
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
	ct.Foreground(ct.Green, true)
	fmt.Println("### QuickJobs")
	ct.Foreground(ct.White, false)
	setup := preSetup()
	redisConnection := database.NewRedisConn(setup.Config)
	check, err := validator.Validate(redisConnection, setup.Config)
	if err != nil {
		log.Fatal(err)
	}
	if check == true {
		fmt.Println("ALL GOOD")
	}
	defer func() {
		fmt.Println(redisConnection.CloseRedisConnection())
	}()
}
