package service

import (
	"context"
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"
	// "reflect"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tools/common/docs"
	"github.com/tools/common/util"
)

var _logger = logrus.New()

// APIResponse returns the service response
type APIResponse struct {
	Message   string      `json:"serviceMessage"`
	Payload   interface{} `json:"payload,omitempty"`
	ServiceTS string      `json:"ts"`
	IsSuccess bool        `json:"isSuccess"`
	Token     *string     `json:"token,omitempty"`
}

// CommonRestService implements the rest API for HRM transactions
type CommonRestService struct {
	dbUtil      *util.PGSqlDBUtil
	authService *AuthService
}

// NewCommonRestService retuens a new initialized version of the service
func NewCommonRestService(config []byte, verbose bool) *CommonRestService {
	service := new(CommonRestService)
	if err := service.Init(config, verbose); err != nil {
		_logger.Errorf("Unable to intialize service instance %v", err)
		return nil
	}
	return service
}

// Init initializes the service
func (srv *CommonRestService) Init(config []byte, verbose bool) error {
	if verbose {
		_logger.SetLevel(logrus.DebugLevel)
	}
	dbUtil, err := util.NewPGSqlDBUtil(config, verbose)
	if err != nil {
		_logger.Errorf("Error in intializing PGSQL util %v", err)
		return err
	}
	srv.dbUtil = dbUtil

	authService := NewAuthService(dbUtil.GetConnection(), verbose)
	if authService == nil {
		_logger.Errorf("Error in intializing AuthService")
		return fmt.Errorf("error in intializing AuthService")
	}
	srv.authService = authService

	_logger.Info("Service instance intialized. Waiting to lauch the service")

	return nil
}

// Serve runs the infinite service method.
func (srv *CommonRestService) Serve(address string, port int, stopSignal chan bool) {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// TODO: Following to be changed for production
	router.MaxMultipartMemory = 8 << 21 // 16 MB Max file size
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowMethods("GET", "POST", "DELETE", "PUT", "OPTIONS", "HEAD")
	corsConfig.AddAllowHeaders("Authorization", "WWW-Authenticate", "Content-Type", "Accept", "X-Requested-With")
	corsConfig.AddExposeHeaders("Authorization", "WWW-Authenticate", "Content-Type", "Accept", "X-Requested-With")
	router.Use(cors.New(corsConfig))
	// JWT Auth handler
	// router.Use(func(c *gin.Context) {
	// 	authorizationRepository := new(repository.AuthorizationRepository)
	// 	validateAuthorizationOutput := authorizationRepository.ValidateAuthorization(c)
	// 	if validateAuthorizationOutput.IsSuccess {
	// 		c.Next()
	// 		return
	// 	}
	// 	c.JSON(http.StatusUnauthorized, validateAuthorizationOutput)
	// 	c.Abort()
	// })
	router.Use(func(c *gin.Context) {
		authorizationRepository := new(Authorization)

		// First, try the first authorization method.
		validateAuthorizationOutputV2 := authorizationRepository.ValidateAuthorization_V2(c)
		if validateAuthorizationOutputV2.IsSuccess {
			c.Next()
			return
		}

		// If the first method fails, try the second authorization method.
		// validateAuthorizationOutput := authorizationRepository.ValidateAuthorization(c)
		// if validateAuthorizationOutput.IsSuccess {
		// 	c.Next()
		// 	return
		// }

		// If both methods fail, return an unauthorized response and abort the request.
		c.JSON(http.StatusUnauthorized, validateAuthorizationOutputV2)
		c.Abort()
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, buildResponse(true, "Service Available", nil))
	})

	srv.authService.AddRouters(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// ReportGeneration service
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: router,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// service connections
		_logger.Infof("Started listening@%v", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			_logger.Errorf("Unable to listen: %v\n", err)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-stopSignal
		_logger.Infof("Received stop signal")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_logger.Infof("Sending shutdown to http server")
		httpServer.Shutdown(ctx)
	}()
	time.Sleep(1 * time.Second)
	// Wait indefinitely
	wg.Wait()
}

func buildResponse(isOk bool, msg string, payload interface{}) APIResponse {
	return APIResponse{
		IsSuccess: isOk,
		Message:   msg,
		Payload:   payload,
		ServiceTS: time.Now().Format("2006-01-02-15:04:05.000"),
	}
}
