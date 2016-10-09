package config

import (
	"errors"

	multierror "github.com/hashicorp/go-multierror"
)

type Config struct {
	Port  int `json:"Port"`
	Redis struct {
		Host     string `json:"Host"`
		Port     int    `json:"Port"`
		Password string `json:"Password"`
		DB       int    `json:"DB"`
	} `json:"Redis"`
	Postgre struct {
		Host         string `json:"Host"`
		Port         int    `json:"Port"`
		DatabaseName string `json:"DatabaseName"`
		Username     string `json:"Username"`
		Password     string `json:"Password"`
	} `json:"Postgre"`
}

func (conf Config) Validate() (bool, error) {
	var result error
	var check bool = true
	if conf.Port <= 0 || conf.Port > 65535 {
		check = false
		result = multierror.Append(result, errors.New("Config Main Port out of Range"))
	}
	if check == true {
		return true, nil
	}
	return false, result
}
