package logging

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(level, format string) *Logger {
	logger := logrus.New()

	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// Set log format
	if format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	logger.SetOutput(os.Stdout)

	return &Logger{logger}
}

// LogMCPRequest logs an MCP requests with contextual information
func (l *Logger) LogMCPRequest(method string, uri string, params interface{}) {
	l.WithFields(logrus.Fields{
		"component": "MCP",
		"method":    method,
		"uri":       uri,
		"params":    params,
	}).Info("Processing MCP request")
}

// LogMCPResponse logs an MCP response with timing information
func (l *Logger) LogMCPResponse(method string, duration time.Duration, err error) {
	fields := logrus.Fields{
		"component": "MCP",
		"method":    method,
		"duration":  duration.String(),
	}

	if err != nil {
		l.WithFields(fields).WithError(err).Error("MCP request failed")
	} else {
		l.WithFields(fields).Info("MCP request completed successfully")
	}
}

// LogK8sOperation logs Kubernetes operations
func (l *Logger) LogK8sOperation(operation string, namespace string, resource string, duration time.Duration, err error) {
	fields := logrus.Fields{
		"component": "Kubernetes",
		"operation": operation,
		"namespace": namespace,
		"resource":  resource,
		"duration":  duration.String(),
	}

	if err != nil {
		l.WithFields(fields).WithError(err).Error("Kubernetes operation failed")
	} else {
		l.WithFields(fields).Info("Kubernetes operation completed successfully")
	}
}
