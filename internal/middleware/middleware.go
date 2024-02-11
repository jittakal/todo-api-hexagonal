package middleware

import (
	"context"
)

type contextKey int

const (
	correlationLogIDKey contextKey = iota
)

func WithCorrelationLogID(ctx context.Context, logID string) context.Context {
	return context.WithValue(ctx, correlationLogIDKey, logID)
}

func GetCorrelationLogID(ctx context.Context) string {
	logID, _ := ctx.Value(correlationLogIDKey).(string)
	return logID
}
