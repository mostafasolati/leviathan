package logger

import (
	"fmt"

	"github.com/mostafasolati/leviathan/contracts"
)

type logger struct {
}

// NewLogger creates a new ILogger.
func NewLogger(configService contracts.IConfigService) contracts.ILogger {
	return &logger{}
}

func (s *logger) Trace(message string) {
	fmt.Println(message)
}

func (s *logger) Debug(message string) {
	fmt.Println(message)
}

func (s *logger) Info(message string) {
	fmt.Println(message)
}

func (s *logger) Warn(message string) {
	fmt.Println(message)
}

func (s *logger) Error(message string) {
	fmt.Println(message)
}

func (s *logger) Fatal(message string) {
	fmt.Println(message)
}

func (s *logger) WithFields(fields contracts.LogFields) contracts.ILogger {
	return s
}
