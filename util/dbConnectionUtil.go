package util

import (
	_ "context"
	"fmt"
	"log"
	"os"
	_ "strconv"

	_ "github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v4/log/logrusadapter"
	_ "github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// CreateConnectionUsingGormToProcurementSchema creates database connection using gorm to procurement schema
func CreateConnectionUsingGormToCommonSchema() *gorm.DB {
	fmt.Println("Connecting....")
	dbHost := viper.GetViper().GetString("db_host")
	dbPort := viper.GetViper().GetString("db_port")
	dbName := viper.GetViper().GetString("db_name")
	dbUsername := viper.GetViper().GetString("db_username")
	dbPassword := viper.GetViper().GetString("db_password")

	dataSourceName := "host=" + dbHost + " user=" + dbUsername + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // Log level
			Colorful: true,
		},
	)
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   viper.GetViper().GetString("common_schema_name") + ".",
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		// fmt.Println("failed to connect database")
		panic(err)
	} else {
		return db
	}
}

// ViperReturnStringConfigVariableFromLocalConfigJSON returns values of string variable from local-config.json
func ViperReturnStringConfigVariableFromLocalConfigJSON(key string) string {
	return viper.GetViper().GetString(key)
}

// ViperReturnIntegerConfigVariableFromLocalConfigJSON returns values of int variable from local-config.json
func ViperReturnIntegerConfigVariableFromLocalConfigJSON(key string) int {
	// viper.SetConfigFile("local-config.json")
	var fileDetails string = ConfigFileName
	// fmt.Println("File Name1 :", fileDetails)
	var (
		fileName string
		fileType string
		location string
	)
	if fileDetails != "" {
		fileName, fileType, location = ReturnConfigFileDetails(fileDetails)
	}

	viper.SetConfigName(fileName) // name of config file (without extension)
	viper.SetConfigType(fileType) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(location) // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value := viper.GetInt(key)
	return value
}
