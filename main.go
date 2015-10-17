package main

import (
	"flag"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/zak01011996/flag-parser/conf"
)

// Main config struct
type CommndLineArguments struct {
	Configfile string `required:"false" name:"config" default:"/etc/daemon.conf" description:"Конфигурационный файл"`
	Daemon     bool   `required:"true" name:"daemon" default:"false" description:"Запуск приложения в режиме daemon"`
	Test       uint32 `required:"false" name:"test" default:"200" description:"Test field"`
}

// Initializing logger
var log = logrus.New()

func main() {
	// Initializing config struct
	config := CommndLineArguments{}

	// Trying to fill config
	if err := conf.GetArguments(&config); err != nil {
		log.Error(err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Checking result
	log.Infof("Config data: %+v", config)
}
