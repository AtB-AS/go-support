package logging

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

// ConfigureGcpLogging configures Zerolog to use field names and values that are understood by GCP
func ConfigureGcpLogging() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.LevelFieldName = "severity"
	zerolog.LevelPanicValue = "critical"
	zerolog.LevelFatalValue = "critical"
}

// AccountIdContext adds an account-based ticketing customer account ID to the logging context.
func AccountIdContext(ctx context.Context, logger zerolog.Logger, accountId string) (context.Context, zerolog.Logger) {
	return addLoggingContext(ctx, logger, "customerAccountId", accountId)
}

// ProfileContext adds a customer profile number to the logging context.
func ProfileContext(ctx context.Context, logger zerolog.Logger, customerNumber string) (context.Context, zerolog.Logger) {
	return addLoggingContext(ctx, logger, "customerNumber", customerNumber)
}

// FirebaseContext adds a Firebase Auth user ID to the logging context.
func FirebaseContext(ctx context.Context, logger zerolog.Logger, uid string) (context.Context, zerolog.Logger) {
	return addLoggingContext(ctx, logger, "firebaseUid", uid)
}

func addLoggingContext(ctx context.Context, logger zerolog.Logger, name string, value string) (context.Context, zerolog.Logger) {
	logger = logger.With().Str(name, value).Logger()
	ctx = logger.WithContext(ctx)

	return ctx, logger
}
