package grpc

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RequestData holds logging data for a gRPC request.
type RequestData struct {
	Method   string
	Request  interface{}
	Metadata metadata.MD
	Latency  time.Duration
	Err      error
}

// MarshalZerologObject formats gRPC request data for Zerolog.
func (c RequestData) MarshalZerologObject(e *zerolog.Event) {
	e.Str("method", c.Method).
		Interface("request", c.Request).
		Interface("metadata", c.Metadata).
		Dur("latencyMs", c.Latency).
		AnErr("err", c.Err)
}

// RequestLoggingInterceptor creates a new interceptor for logging gRPC requests.
func RequestLoggingInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		reqData := RequestData{
			Method:  method,
			Request: req,
		}

		// Add metadata if present
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			reqData.Metadata = md
		}

		// Invoke call
		begin := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Log result
		reqData.Latency = time.Since(begin)
		if err == nil {
			log.Ctx(ctx).Info().Interface("grpc", reqData).Msgf("Completed call to %s", method)
		} else {
			reqData.Err = err
			log.Ctx(ctx).Info().Interface("grpc", reqData).Msgf("Completed call to %s with error", method)
		}

		return err
	}
}
