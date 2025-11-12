package service

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/rest/api/internal/model"
	"github.com/rest/api/internal/util"

	"github.com/sirupsen/logrus"
)

var _asLogger = logrus.New()

// RESTService provides authentication related rest services
type RESTService struct {
	dbConn        *util.DBConnectionWrapper
	jwtSigningKey []byte
	bypassAuth    map[string]bool
}

// NewAuthenticationRESTService returns a new initialized version of the service
func NewAuthenticationRESTService(config []byte, dbConnection *util.DBConnectionWrapper, verbose bool) *RESTService {
	service := new(RESTService)
	if err := service.Init(config, dbConnection, verbose); err != nil {
		_asLogger.Errorf("Unable to initialize service instance %v", err)
		return nil
	}
	return service
}

// Init initializes the service instance
func (s *RESTService) Init(config []byte, dbConnection *util.DBConnectionWrapper, verbose bool) error {
	if verbose {
		_asLogger.SetLevel(logrus.DebugLevel)
	}
	if dbConnection == nil {
		return fmt.Errorf("null DB Util reference passed")
	}
	s.dbConn = dbConnection
	var conf model.AuthServiceConfig
	err := json.Unmarshal(config, &conf)
	if err != nil {
		_asLogger.Error("Unable to parse config json file ", err)
		return err
	}
	if conf.JWTKey != nil && len(*conf.JWTKey) > 0 {
		s.jwtSigningKey = []byte(*conf.JWTKey)
	}
	s.bypassAuth = make(map[string]bool)
	s.bypassAuth["/"] = true
	if conf.BypassAuth != nil && len(conf.BypassAuth) > 0 {
		for _, url := range conf.BypassAuth {
			s.bypassAuth[url] = true
		}
	}
	_asLogger.Infof("Successfully initialized AuthenticationRESTService")
	return nil
}

// AddRouters add api end points specific to this service
func (s *RESTService) AddRouters(apiBase string, router *gin.Engine) {
	router.POST("/api/auth/create", func(c *gin.Context) {
		resp := s.createUser(c)
		c.JSON(resp.StatusCode, resp)
	})

	router.POST("/api/auth/login", func(c *gin.Context) {
		resp := s.validateLogin(c)
		c.JSON(resp.StatusCode, resp)
	})

	router.POST("/api/auth/resetpwd", func(c *gin.Context) {
		resp := s.resetPassword(c)
		c.JSON(resp.StatusCode, resp)
	})

	router.PUT("/api/auth/update", func(c *gin.Context) {
		resp := s.updateUser(c)
		c.JSON(resp.StatusCode, resp)
	})

	router.GET("/api/auth/users", func(c *gin.Context) {
		resp := s.getAllUsers(c)
		c.JSON(resp.StatusCode, resp)
	})

	// Satcom Data CRUD routes
	router.POST("/api/satcom", func(c *gin.Context) {
		resp := s.createSatcomData(c)
		c.JSON(resp.StatusCode, resp)
	})

	router.GET("/api/satcom", func(c *gin.Context) {
		resp := s.getAllSatcomData(c)
		c.JSON(resp.StatusCode, resp)
	})

	router.GET("/api/satcom/:id", func(c *gin.Context) {
		resp := s.getSatcomDataById(c)
		c.JSON(resp.StatusCode, resp)
	})

	router.PUT("/api/satcom/:id", func(c *gin.Context) {
		resp := s.updateSatcomData(c)
		c.JSON(resp.StatusCode, resp)
	})

	router.DELETE("/api/satcom/:id", func(c *gin.Context) {
		resp := s.deleteSatcomData(c)
		c.JSON(resp.StatusCode, resp)
	})
}
