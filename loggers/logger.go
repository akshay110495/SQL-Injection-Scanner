package loggers

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func GetLogger() *log.Logger {
	logMap := log.FieldMap{
		log.FieldKeyTime: "@timestamp",
		log.FieldKeyMsg:  "message",
	}

	var logger = log.New()
	logger.SetFormatter(&log.JSONFormatter{FieldMap: logMap})
	logger.SetOutput(io.MultiWriter(os.Stdout))
	return logger
}

func WriteLog(logger *log.Logger, err error) {
	logger.WithFields(log.Fields{}).Info(err)
}
