package main

import (
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rdoorn/gohelper/app"
)

type MyApp struct {
	//app.Webserver
	app.App
}

type MyAppConfig struct {
	app.App
	Webserver app.WebserverConfig
	//Profiling profiling.Config
}

func main() {
	router := gin.Default()

	handler := &MyApp{}
	handler.New(
		app.Webserver(app.WebserverConfig{
			IP:        "localhost",
			Port:      8080,
			TLSConfig: nil,
			Handler:   router,
		}),
		app.Signal(handler.Stop, os.Interrupt, syscall.SIGTERM),
		//app.Signal(handler.reload, syscall.SIGHUP),
	)

	// all requests require authentication using certificates
	//router.Use(handler.AuthClientCertificate)
	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/version", handler.version)
	}

	/*app.Signal(func() {
		app.Stop()
	}, os.Interrupt, syscall.SIGTERM)

	app.Signal(func() {
		app.Reload()
	}, syscall.SIGHUP)
	*/

	if err := handler.Start(); err != nil {
		handler.Panicf("failed to start server: %s", err)
	}

}

func (m *MyApp) stop() {
	m.Stop()
}

/*
func (m *MyApp) reload() {
	m.Reload()
}
*/

func (m *MyApp) version(c *gin.Context) {
}
