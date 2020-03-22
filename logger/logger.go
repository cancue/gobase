package logger

import (
	"fmt"

	"github.com/bluele/logrus_slack"
	"github.com/sirupsen/logrus"

	"github.com/cancue/gobase/config"
	"github.com/cancue/gobase/errors"
)

var logger Logger

// Logger is exported.
type Logger interface {
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
}

// DefaultLogger is exported.
type DefaultLogger struct {
	*logrus.Entry
}

// Error public
func (lgr *DefaultLogger) Error(args ...interface{}) {
	err := args[0]

	eErr, ok := err.(*errors.Err)
	if !ok {
		lgr.Entry.Error(err)
		return
	}

	lgr.Entry.
		WithFields(logrus.Fields{
			"trace": fmt.Sprintf("%+v", eErr.Trace()),
			"data":  fmt.Sprintf("%+v", eErr.Data),
		}).
		Error(eErr)

	return
}

// Get retrives logger
func Get() Logger {
	return logger
}

// Set public
func Set(conf *config.Config) (defaultLogger *DefaultLogger) {
	lgr := logrus.New()
	lgr.SetLevel(logrus.DebugLevel)

	if slack, ok := conf.YAML["slack"].(map[string]string); ok {
		var level logrus.Level

		switch slack["accepted-level"] {
		case "panic":
			level = logrus.PanicLevel
		case "fatal":
			level = logrus.FatalLevel
		case "error":
			level = logrus.ErrorLevel
		case "warn":
			level = logrus.WarnLevel
		case "info":
			level = logrus.InfoLevel
		case "debug":
			level = logrus.DebugLevel
		case "trace":
			level = logrus.TraceLevel
		default:
			level = logrus.WarnLevel
		}

		lgr.AddHook(&logrus_slack.SlackHook{
			HookURL:        slack["hook-url"],
			AcceptedLevels: logrus_slack.LevelThreshold(level),
			Channel:        slack["channel"],
			IconEmoji:      slack["icon"],
			Username:       slack["username"],
		})
	}

	defaultLogger = &DefaultLogger{
		lgr.WithFields(
			logrus.Fields{
				"service": conf.YAML["name"].(string),
				"stage":   conf.Stage,
			},
		),
	}
	logger = defaultLogger

	return
}

// SetCustomLogger is exported.
func SetCustomLogger(lgr Logger) {
	logger = lgr
}
