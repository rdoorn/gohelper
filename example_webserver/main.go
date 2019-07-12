package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rdoorn/gohelper/app"
	"github.com/rdoorn/gohelper/cmd"
	"github.com/spf13/cobra"
)

type MyApp struct {
	//app.Webserver
	app.App
}

type MyAppConfig struct {
	app.App
	Webserver app.WebserverConfig
}

func server() func(cmd *cobra.Command, args []string) {
	return func(cobra *cobra.Command, args []string) {
		router := gin.New()

		handler := &MyApp{}
		handler.New(
			app.Webserver(app.WebserverConfig{
				IP:        cmd.MustGetString(cobra, "listener"),
				Port:      cmd.MustGetInt(cobra, "port"),
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
			handler.Panicf(err)
		}
	}
}

func main() {
	app.ServerCmd.Run = server()

	if err := app.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}

func init() {
	log.Println("init app")

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
