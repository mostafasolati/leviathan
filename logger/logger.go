package logger

import (
	"log/syslog"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	rollingHook "github.com/lanziliang/logrus-rollingfile-hook"
	"github.com/mostafasolati/leviathan/contracts"
	"github.com/sirupsen/logrus"
	syslogHook "github.com/sirupsen/logrus/hooks/syslog"
	elogrus "gopkg.in/go-extras/elogrus.v7"
)

const (
	// CPLogLevel is the configuration for the log level. Available options are
	// trace, debug, info, warn, error and fatal (default: info).
	CPLogLevel = "log.level"

	// CPLogAppender is the configuration for the medium on-which logs are
	// written. Available options are stderr, file and syslog (default: stderr).
	CPLogAppender = "log.appender"

	// CPLogFile is the configuration for the rolling log file's path, e.g.
	CPLogFile = "log.file"

	// CPLogSyslogHost is the configuration for syslog host and port, e.g.
	// localhost:514.
	CPLogSyslogHost = "log.syslog.host"

	// CPLogElasticAddress is the configuration for elasticsearch endpoits, e.g.
	// http://localhost:9001.
	CPLogElasticAddress = "log.elastic.address"

	// CPLogElasticUsername is the configuration for elasticsearch username.
	CPLogElasticUsername = "log.elastic.username"

	// CPLogElasticPassword is the configuration for elasticsearch password.
	CPLogElasticPassword = "log.elastic.password"

	// CPLogElasticIndex is the configuration for elasticsearch log index.
	CPLogElasticIndex = "log.elastic.index"
)

var logLevelNames = map[string]logrus.Level{
	"trace": logrus.TraceLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"fatal": logrus.FatalLevel,
}

type logBackendFactory = func(contracts.IConfigService) (*logrus.Logger, error)

var logBackendFactories = map[string]logBackendFactory{
	"stderr":  createStderrLogBackend,
	"file":    createFileLogBackend,
	"syslog":  createSyslogLogBackend,
	"elastic": createElasticLogBackend,
}

type logger struct {
	backend *logrus.Logger
}

// NewLogger creates a new ILogger.
func NewLogger(configService contracts.IConfigService) (contracts.ILogger, error) {
	appender := configService.String(CPLogAppender)
	factory, ok := logBackendFactories[appender]
	if !ok {
		factory = logBackendFactories["stderr"]
	}
	backend, err := factory(configService)
	if err != nil {
		return nil, err
	}
	backend.Level = logLevelFromConfig(configService)
	return &logger{backend: backend}, nil
}

func logLevelFromConfig(config contracts.IConfigService) logrus.Level {
	levelName := config.String(CPLogLevel)
	if level, ok := logLevelNames[levelName]; ok {
		return level
	}
	return logrus.InfoLevel
}

func createStderrLogBackend(contracts.IConfigService) (*logrus.Logger, error) {
	backend := logrus.New()
	backend.Formatter = &logrus.TextFormatter{
		DisableSorting:   true,
		PadLevelText:     true,
		QuoteEmptyFields: true,
	}
	return backend, nil
}

func createFileLogBackend(config contracts.IConfigService) (*logrus.Logger, error) {
	logFile := config.String(CPLogFile)
	hook, err := rollingHook.NewRollingFileTimeHook(logFile, "2006-01-02", 0)
	if err != nil {
		return nil, err
	}
	backend, _ := createStderrLogBackend(config)
	backend.Hooks.Add(hook)
	return backend, nil
}

func createSyslogLogBackend(config contracts.IConfigService) (*logrus.Logger, error) {
	host := config.String(CPLogSyslogHost)
	hook, err := syslogHook.NewSyslogHook("udp", host, syslog.LOG_INFO, "")
	if err != nil {
		return nil, err
	}
	backend, _ := createStderrLogBackend(config)
	backend.Hooks.Add(hook)
	return backend, nil
}

func createElasticLogBackend(config contracts.IConfigService) (*logrus.Logger, error) {
	addressConfig := config.String(CPLogElasticAddress)
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: strings.Split(addressConfig, ","),
		Username:  config.String(CPLogElasticUsername),
		Password:  config.String(CPLogElasticPassword),
	})
	if err != nil {
		return nil, err
	}
	hook, err := elogrus.NewAsyncElasticHook(
		client,
		config.BaseURL(),
		logLevelFromConfig(config),
		config.String(CPLogElasticIndex),
	)
	if err != nil {
		return nil, err
	}
	backend, _ := createStderrLogBackend(config)
	backend.Hooks.Add(hook)
	return backend, nil
}

// Trace implements ILogger.Trace
func (log *logger) Trace(message string) {
	log.backend.Trace(message)
}

// Debug implements ILogger.Debug
func (log *logger) Debug(message string) {
	log.backend.Debug(message)
}

// Info implements ILogger.Info
func (log *logger) Info(message string) {
	log.backend.Info(message)
}

// Warn implements ILogger.Warn
func (log *logger) Warn(message string) {
	log.backend.Warn(message)
}

// Error implements ILogger.Error
func (log *logger) Error(message string) {
	log.backend.Error(message)
}

// Fatal implements ILogger.Fatal
func (log *logger) Fatal(message string) {
	log.backend.Fatal(message)
}

// WithFields implements ILogger.WithFields
func (log *logger) WithFields(fields contracts.LogFields) contracts.ILogger {
	return &loggerWithFields{
		backend: log.backend,
		fields:  fields,
	}
}

type loggerWithFields struct {
	backend *logrus.Logger
	fields  contracts.LogFields
}

// Trace implements ILogger.Trace
func (log *loggerWithFields) Trace(message string) {
	log.backend.WithFields(logrus.Fields(log.fields)).Trace(message)
}

// Debug implements ILogger.Debug
func (log *loggerWithFields) Debug(message string) {
	log.backend.WithFields(logrus.Fields(log.fields)).Debug(message)
}

// Info implements ILogger.Info
func (log *loggerWithFields) Info(message string) {
	log.backend.WithFields(logrus.Fields(log.fields)).Info(message)
}

// Warn implements ILogger.Warn
func (log *loggerWithFields) Warn(message string) {
	log.backend.WithFields(logrus.Fields(log.fields)).Warn(message)
}

// Error implements ILogger.Error
func (log *loggerWithFields) Error(message string) {
	log.backend.WithFields(logrus.Fields(log.fields)).Error(message)
}

// Fatal implements ILogger.Fatal
func (log *loggerWithFields) Fatal(message string) {
	log.backend.WithFields(logrus.Fields(log.fields)).Fatal(message)
}

// WithFields implements ILogger.WithFields
func (log *loggerWithFields) WithFields(fields contracts.LogFields) contracts.ILogger {
	for field, value := range fields {
		log.fields[field] = value
	}
	return log
}
