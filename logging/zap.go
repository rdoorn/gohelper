package logging

import (
	"log/syslog"

	"github.com/tchap/zapext/zapsyslog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapSugar struct {
	Logger *zap.SugaredLogger
	sync   func() error
}

func (f *ZapSugar) Println(v ...interface{}) {
	switch v[0].(type) {
	case string:
		f.Logger.Infow(v[0].(string), v[1:]...)
	case error:
		f.Logger.Infow(v[0].(error).Error(), v[1:]...)
	default:
		f.Logger.Infow("Info", v...)
	}
}

func (f *ZapSugar) Debugf(v ...interface{}) {
	switch v[0].(type) {
	case string:
		f.Logger.Debugw(v[0].(string), v[1:]...)
	case error:
		f.Logger.Debugw(v[0].(error).Error(), v[1:]...)
	default:
		f.Logger.Debugw("Debug", v...)
	}
}

func (f *ZapSugar) Infof(v ...interface{}) {
	switch v[0].(type) {
	case string:
		f.Logger.Infow(v[0].(string), v[1:]...)
	case error:
		f.Logger.Infow(v[0].(error).Error(), v[1:]...)
	default:
		f.Logger.Infow("Info", v...)
	}
}

func (f *ZapSugar) Warnf(v ...interface{}) {
	switch v[0].(type) {
	case string:
		f.Logger.Warnw(v[0].(string), v[1:]...)
	case error:
		f.Logger.Warnw(v[0].(error).Error(), v[1:]...)
	default:
		f.Logger.Warnw("Warn", v...)
	}
}

func (f *ZapSugar) Errorf(v ...interface{}) {
	switch v[0].(type) {
	case string:
		f.Logger.Errorw(v[0].(string), v[1:]...)
	case error:
		f.Logger.Errorw(v[0].(error).Error(), v[1:]...)
	default:
		f.Logger.Errorw("Error", v...)
	}
}

func (f *ZapSugar) Fatalf(v ...interface{}) {
	switch v[0].(type) {
	case string:
		f.Logger.Fatalw(v[0].(string), v[1:]...)
	case error:
		f.Logger.Fatalw(v[0].(error).Error(), v[1:]...)
	default:
		f.Logger.Fatalw("Fatal", v...)
	}
}

func (f *ZapSugar) Panicf(v ...interface{}) {
	switch v[0].(type) {
	case string:
		f.Logger.Panicw(v[0].(string), v[1:]...)
	case error:
		f.Logger.Panicw(v[0].(error).Error(), v[1:]...)
	default:
		f.Logger.Panicw("Panic", v...)
	}
}

/*
func (f *ZapSugar) Loglevel(level string) {
	switch level {
	case "debug":
		f.Logger.With(zap.Le .Level(zap.DebugLevel)

	}
}
*/

func NewZapSyslog(dst ...string) (*ZapSugar, error) {
	// Initialize a syslog writer.
	writer, err := syslog.New(syslog.LOG_ERR|syslog.LOG_LOCAL0, "")
	if err != nil {
		return nil, err
	}

	// Initialize Zap.
	encoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapsyslog.NewCore(zapcore.ErrorLevel, encoder, writer)

	logger := zap.New(core, zap.Development(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &ZapSugar{
		Logger: logger.Sugar(),
		sync:   logger.Sync,
	}, err

}

func NewZapConsole(dst ...string) (*ZapSugar, error) {

	//config := zap.New(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()))
	//config := zap.New(zapcore.NewCore(zapcore.NewTextEncoder(zapcore.TextTimeFormat(time.RFC822), zap.NewProductionConfig())))

	/*zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(NewProductionConfig().EncoderConfig),
			zapcore.Lock(zapcore.AddSync(ztest.FailWriter{})),
			DebugLevel,
		),
		ErrorOutput(errSink),
	)*/

	config := zapConfig(dst...)
	config.Encoding = "console"
	//config.EncoderConfig = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	//config.EncoderConfig = )
	logger, err := config.Build()

	return &ZapSugar{
		Logger: logger.Sugar(),
		sync:   logger.Sync,
	}, err
}

func NewZap(dst ...string) (*ZapSugar, error) {
	config := zapConfig(dst...)
	logger, err := config.Build()

	return &ZapSugar{
		Logger: logger.Sugar(),
		sync:   logger.Sync,
	}, err
}

func zapConfig(dst ...string) zap.Config {
	config := zap.NewProductionConfig()
	if len(dst) == 0 {
		config.OutputPaths = []string{"stdout"}
	} else {
		config.OutputPaths = dst
	}
	config.DisableCaller = true
	// config.Encoding = "console"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.TimeKey = "timestamp"
	//config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.MessageKey = "message"
	config.Level.SetLevel(zap.DebugLevel)
	//cfg.Level = zap.DebugLevel
	return config
}
