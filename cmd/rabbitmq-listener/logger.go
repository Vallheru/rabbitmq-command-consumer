package main

import (
	"go.uber.org/zap"
)

// The Logger interface is a thin layer to provide simple
// interface for logging in the system
type Logger interface {
	Info(mod string, msg string, keysAndValues ...interface{})
	Error(mod string, msg string, keysAndValues ...interface{})
	Warn(mod string, msg string, keysAndValues ...interface{})
	Destroy()
}

// NewZapLogger ...
func NewZapLogger(path string) (*ZapLogger, error) {
	cfg := zap.NewProductionConfig()
	
	if path != "" {
		cfg.OutputPaths = []string{
			path,
		}
	}
	logger, err := cfg.Build()

	if err != nil {
		return nil, err
	}

	return &ZapLogger{log: logger}, nil
}

// ZapLogger is logging wrapper for Uber Zap logging library
type ZapLogger struct {
	Logger
	
	log *zap.Logger
}

// Info is function to log message with the "info" level
func (l *ZapLogger) Info(mod string, msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Infow("[ " + mod + " ] " + msg, 
		keysAndValues... ,
	)
}

// Error is function to log message with the "error" level
func (l *ZapLogger) Error(mod string, msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Errorw("[ " + mod + " ] " + msg, 
		keysAndValues... ,
	)
}

// Warn is function to log message with the "warn" level
func (l *ZapLogger) Warn(mod string, msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Warnw("[ " + mod + " ] " + msg, 
		keysAndValues... ,
	)
}

// Destroy is function clear memory from unused stuff
func (l *ZapLogger) Destroy() {
	l.log.Sync()
}

// LoggerMock ...
type LoggerMock struct {
	Logger
	output string
}

// Info ...
func (*LoggerMock) Info(mod string, msg string, keysAndValues ...interface{}) {
}

// Error ...
func (*LoggerMock) Error(mod string, msg string, keysAndValues ...interface{}) {
}

// Warn ...
func (*LoggerMock) Warn(mod string, msg string, keysAndValues ...interface{}) {
}

// Destroy ...
func (*LoggerMock) Destroy() {
}
