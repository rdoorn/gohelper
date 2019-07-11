package main

import (
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rdoorn/gohelper/app"
)

type MyApp struct {
	app.Webserver
}

func main() {
	webserver := &MyApp{}
	webserver.Init()

	/*myApp, err := app.New(webserver)
	if err != nil {
		myApp.Panicf("failed to initialize webserver: %s", err)
	}*/

	router := gin.Default()
	// all requests require authentication using certificates
	//router.Use(handler.AuthClientCertificate)
	// Simple group: v1
	v1 := router.Group("/v3")
	{
		v1.POST("/search/product", handler.ProductRequest)
	}

	webserver.LoadWebserverConfig(app.WebserverConfig{
		IP:        "localhost",
		Port:      8080,
		TLSConfig: nil,
		Handler:   router,
	})

	webserver.Signal(func() {
		webserver.Stop()
	}, os.Interrupt, syscall.SIGTERM)

	webserver.Signal(func() {
		webserver.Reload()
	}, syscall.SIGHUP)

	if err := webserver.Start(); err != nil {
		webserver.Panicf("failed to start webserver: %s", err)
	}

}
