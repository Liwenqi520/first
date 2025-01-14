package log

import "context"

type userKey int
type traceKey int

var traceIDKey traceKey
var userIDKey userKey

func NewTraceIDContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey, id)
}

func FromTraceIDContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(traceIDKey).(string)
	return id, ok
}

func NewUserIDContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func FromUserIDContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
