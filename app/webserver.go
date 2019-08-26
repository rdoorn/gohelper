package app

/*
type WebserverConfig struct {
	IP        string
	Port      int
	Handler   http.Handler
	TLSConfig *tls.Config
}

type WebserverHandler struct {
	logging    logging.SimpleLogger
	config     *WebserverConfig
	server     http.Server
	serverDone chan struct{}
}

var ServerCmd = &cobra.Command{
	Short:  "Server starts the webserver",
	Use:    "server",
	PreRun: Lock(),
}

func init() {
	log.Println("init webserver")
	ServerCmd.PersistentFlags().Bool("skip-tls-verify", false, "Foolishly accept TLS certificates signed by unkown certificate authorities")
	ServerCmd.PersistentFlags().String("listener", "127.0.0.1", "Listening address")
	ServerCmd.PersistentFlags().Int("port", 8080, "Listening port")
	ServerCmd.MarkFlagRequired("listener")
	ServerCmd.MarkFlagRequired("port")

	cmd.Root.AddCommand(ServerCmd)
}

func NewWebserverHandler(logger logging.SimpleLogger, c WebserverConfig) *WebserverHandler {
	return &WebserverHandler{
		logging:    logger, // inherrit logger
		config:     &c,
		serverDone: make(chan struct{}),
	}
}

func (w *WebserverHandler) Start() error {
	// router.POST("/opvragen/naw", handler.ClientRequest)
	// router.POST("/1bnr", handler.IbanRequest)

	listenAddr := fmt.Sprintf("%s:%d", w.config.IP, w.config.Port)
	// start https server
	w.server = http.Server{
		Addr:      listenAddr,
		Handler:   w.config.Handler,
		TLSConfig: w.config.TLSConfig,
	}

	w.logging.Infof("Webservice listening", "addr", listenAddr)
	if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	<-w.serverDone
	return nil
}

func (w *WebserverHandler) Stop() error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	w.server.SetKeepAlivesEnabled(false)
	if err := w.server.Shutdown(ctx); err != nil {
		w.logging.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(w.serverDone)
	w.logging.Infof("Webservice shutdown")

	return nil
}

func (w *WebserverHandler) Reload() error {
	return nil
}

func (a *App) GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		end := time.Now()
		if raw != "" {
			path = path + "?" + raw
		}

		a.Infof(
			//path,
			c.Request.Host,
			"client", c.ClientIP(),
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency", end.Sub(start),
			"size", c.Writer.Size(),
			"error", c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)

	}
}
*/

/*
func (w *WebserverHandler) LoadWebserverConfig(config WebserverConfig) error {
	w.Config = &config
	return nil
}
*/
/*
func (w *WebserverHandler) Signal(f func(), sig ...os.Signal) {
	w.signal.Add(f, sig...)
}
*/
