package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/viper"
	"github.com/tools/common/service"
	"github.com/tools/common/util"
)

var _version = "default"

func main() {
	configFilePath := flag.String("c", "./config.json", "Configuration file")
	verbose := flag.Bool("v", false, "Verbose")
	port := flag.Int("port", 7070, "Port to run the server on")

	flag.Parse()
	util.ConfigFileName = *configFilePath
	configBytes, err := os.ReadFile(*configFilePath)
	if err != nil {
		fmt.Println("Unable to read configuration file ", err.Error())
		os.Exit(1)
	}
	initAppConfigViper(*configFilePath)
	bindAddress := flag.String("b", "0.0.0.0", "Bind address")
	if commonService := service.NewCommonRestService(configBytes, *verbose); commonService != nil {
		stopSignal := make(chan bool)
		termination := make(chan os.Signal)
		signal.Notify(termination, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-termination
			fmt.Println("SIGTERM/SIGINT received from os")
			stopSignal <- true
		}()
		commonService.Serve(*bindAddress, *port, stopSignal)
	} else {
		fmt.Println("Unable to start the service ...")
		os.Exit(2)
	}
}

func initAppConfigViper(configPath string) {
	viper.SetConfigFile(configPath)
	//viper.SetConfigName("appConfig")
	//viper.SetConfigType("json")
	//viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Panic("There is no such a config file in path ./config/appConfig.json")
		} else {
			log.Panic("There is some problem about data in file")
		}
	}
	viper.AutomaticEnv()
	fmt.Println("Viper automatic env is set: ", viper.GetViper().GetString("db_host"))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
