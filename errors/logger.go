package errors

import (
	"fmt"

	"github.com/bluele/logrus_slack"
	"github.com/sirupsen/logrus"

	"github.com/cancue/gobase/config"
)

// SetLogger set logger in dependency injection.
func SetLogger(conf *config.Config) {
	lgr := logrus.New()
	lgr.SetLevel(logrus.DebugLevel)
	setSlack(lgr, conf)

	logger = lgr.WithFields(
		logrus.Fields{
			"service": conf.Name,
			"stage":   conf.Stage,
		},
	)
}

// LogError .
func LogError(err interface{}) {
	eErr, ok := err.(*richError)
	if !ok {
		logger.Error(err)
		return
	}

	logger.
		WithFields(logrus.Fields{
			"trace": fmt.Sprintf("%+v", eErr.trace()),
			"data":  fmt.Sprintf("%+v", eErr.Data),
		}).
		Error(eErr)

	return
}

/* private */

var logger *logrus.Entry

func setSlack(lgr *logrus.Logger, conf *config.Config) {
	slack, ok := conf.YAML["slack"].(map[string]string)

	if !ok {
		return
	}

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
