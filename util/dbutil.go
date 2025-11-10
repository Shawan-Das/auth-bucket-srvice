package util

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

var _dbUtillog = logrus.New()

// PGSqlDBUtil type implements a postgressql db driver wrapper
type PGSqlDBUtil struct {
	connectionStr string
	dbConnection  *pgx.Conn
	verbose       bool
	mutex         *sync.Mutex
}

type dbConfig struct {
	DBHost   string `json:"dbhost"`
	DBName   string `json:"dnname"`
	UserID   string `json:"uid"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout"`
	Port     int    `json:"port"`
}

// NewPGSqlDBUtil returns a new instance of PGSQLUtil
func NewPGSqlDBUtil(configBytes []byte, verbose bool) (*PGSqlDBUtil, error) {
	_dbUtillog.Info("Creating DBUtil...")
	dbUtil := new(PGSqlDBUtil)
	if err := dbUtil.Init(configBytes); err != nil {
		return nil, err
	}
	if verbose {
		_dbUtillog.SetLevel(logrus.DebugLevel)
	}
	return dbUtil, nil
}

// Init intializes the connection to db
func (dbu *PGSqlDBUtil) Init(configBytes []byte) error {
	var config dbConfig
	err := json.Unmarshal(configBytes, &config)
	if err != nil {
		_dbUtillog.Errorf("Unable to parse the configuration %v", err)
		return err
	}
	dbu.connectionStr = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", config.UserID, config.Password, config.DBHost, config.Port, config.DBName)
	dbu.dbConnection, err = pgx.Connect(context.Background(), dbu.connectionStr)
	if err != nil {
		_dbUtillog.Errorf("Unable to connection to database: %v\n", err)
		return err
	}
	dbu.mutex = new(sync.Mutex)
	dbu.verbose = true
	go dbu.keepAlive()
	return nil
}

func (dbu *PGSqlDBUtil) keepAlive() {
	for {
		<-time.After(5 * time.Minute)
		isSuccess := dbu.refreshConnection()
		_dbUtillog.Infof("Refresh connection %v", isSuccess)
	}
}

func (dbu *PGSqlDBUtil) refreshConnection() bool {
	dbu.mutex.Lock()
	defer dbu.mutex.Unlock()
	if dbu.dbConnection.Ping(context.Background()) != nil {
		_dbUtillog.Infof("Reconnecting....")

		dbConnection, err := pgx.Connect(context.Background(), dbu.connectionStr)
		if err != nil {
			_dbUtillog.Errorf("Unable to connection to database: %v\n", err)
			return false
		}
		dbu.dbConnection = dbConnection
	}
	return true
}

// Query runs a query in the db and reurns pgx.Rows
func (dbu *PGSqlDBUtil) Query(sql string, params ...interface{}) (pgx.Rows, error) {
	dbu.refreshConnection()
	results, err := dbu.dbConnection.Query(context.Background(), sql, params...)
	if err != nil {
		_dbUtillog.Errorf("Error in executing the query %v", err)
		return nil, err
	}
	_dbUtillog.Infof("Execution result %v", results)
	return results, nil
}

// Shutdown close the db connection
func (dbu *PGSqlDBUtil) Shutdown() {
	if dbu.dbConnection != nil {
		dbu.dbConnection.Close(context.Background())
	}
}


