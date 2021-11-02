package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	grpcmeta "google.golang.org/grpc/metadata"
)

// RequestLoggingMiddleware creates a new HTTP interceptor that sets any of correlation ID, request ID and
// install ID that are present as headers in the incoming request as context log fields.
// It will also set correlation ID for outgoing gRPC requests using correlation ID if present,
// else request ID, or a new UUID if none of these are set.
func RequestLoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := log.Logger

		// Get tracing information from request
		correlationId := r.Header.Get("Correlation-Id")
		if correlationId != "" {
			logger = logger.With().Str("correlationId", correlationId).Logger()
		}

		requestId := r.Header.Get("Request-Id")
		if requestId == "" {
			requestId = r.Header.Get("Atb-Request-Id")
		}
		if correlationId != "" {
			logger = logger.With().Str("requestId", requestId).Logger()
		}

		installId := r.Header.Get("Install-Id")
		if installId == "" {
			installId = r.Header.Get("Atb-Install-Id")
		}
		if installId != "" {
			logger = logger.With().Str("installId", requestId).Logger()
		}

		// Set outgoing tracing ID for gRPC
		var outgoingCorrelationId string
		if correlationId != "" {
			outgoingCorrelationId = correlationId
		} else if requestId != "" {
			outgoingCorrelationId = requestId
		} else {
			outgoingCorrelationId = uuid.New().String()
			logger.Debug().Msgf("No tracing ID found in incoming request, using %s as correlation ID", outgoingCorrelationId)
		}
		ctx := grpcmeta.AppendToOutgoingContext(r.Context(), "x-correlation-id", outgoingCorrelationId)

		// Attach logger to request
		logCtx := logger.WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(logCtx))
	}

	return http.HandlerFunc(fn)
}
