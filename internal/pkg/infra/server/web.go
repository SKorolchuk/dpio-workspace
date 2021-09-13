package server

import (
	"context"
	"errors"
	"time"
)

type contextKey int

const key contextKey = 19191

// RequestInfo contains information about request.
type RequestInfo struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// GetRequestInfo returns request information of RequestInfo type.
func GetRequestInfo(ctx context.Context) (*RequestInfo, error) {
	info, found := ctx.Value(key).(*RequestInfo)
	if !found {
		return nil, errors.New("info missing from context")
	}
	return info, nil
}

// GetTraceID returns the trace id from the context.
func GetTraceID(ctx context.Context) string {
	info, found := ctx.Value(key).(*RequestInfo)
	if !found {
		return "00000000-0000-0000-0000-000000000000"
	}
	return info.TraceID
}

// SetStatusCode sets the status code back into the context.
func SetStatusCode(ctx context.Context, statusCode int) error {
	info, found := ctx.Value(key).(*RequestInfo)
	if !found {
		return errors.New("web value missing from context")
	}
	info.StatusCode = statusCode
	return nil
}
