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
}

func main() {
	router := gin.New()

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
	router.Use(gin.Recovery(), handler.GinLogger())
	v1 := router.Group("/v1")
	{
		v1.POST("/version", handler.version)
		v1.GET("/hello/:name", handler.hello)
	}

	if err := handler.Start(); err != nil {
		handler.Panicf("failed to start server: %s", err)
	}

}

/*
func (m *MyApp) reload() {
	m.Reload()
}
*/

func (m *MyApp) version(c *gin.Context) {
	//m.Infof("hello verion request")
}

func (m *MyApp) hello(c *gin.Context) {
	c.JSON(200, gin.H{"name": c.Param("name")})
	//m.Infof("hello %s", c.Param("name"))
}
