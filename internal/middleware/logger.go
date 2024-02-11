package middleware

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, format string, args ...interface{})
}

func NewLogger() Logger {
	return &logger{}
}

func GenerateLogID() string {
	return uuid.New().String()
}

type logger struct{}

func (l *logger) log(level logrus.Level, ctx context.Context, format string, args ...interface{}) {
	logID := GetCorrelationLogID(ctx)
	now := time.Now().Format(time.RFC3339)
	msg := fmt.Sprintf(format, args...)

	_, filename, _, _ := runtime.Caller(2)
	packageName := filepath.Base(filepath.Dir(filename))

	logrus.WithFields(logrus.Fields{
		"timestamp":   now,
		"logID":       logID,
		"packageName": packageName,
		"filename":    filepath.Base(filename),
		"stdout":      true,
	}).Log(level, msg)
}

func (l *logger) Info(ctx context.Context, format string, args ...interface{}) {
	l.log(logrus.InfoLevel, ctx, format, args...)
}

func (l *logger) Error(ctx context.Context, format string, args ...interface{}) {
	l.log(logrus.ErrorLevel, ctx, format, args...)
}
