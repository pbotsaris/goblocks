package socketmode

import (
	"errors"
	"net/http"
	"strings"
)

// Error sentinel values.
var (
	ErrHelloTimeout      = errors.New("timeout waiting for hello message")
	ErrConnectionClosed  = errors.New("connection closed")
	ErrWriteTimeout      = errors.New("write timeout")
	ErrShuttingDown      = errors.New("client is shutting down")
	ErrHandlerTimeout    = errors.New("handler timeout")
	ErrConcurrencyLimit  = errors.New("concurrency limit reached")
)

// PermanentError represents an error that should not be retried.
// Examples: invalid auth token, revoked app, app not installed.
type PermanentError struct {
	Err     error
	Message string
}

func (e *PermanentError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *PermanentError) Unwrap() error {
	return e.Err
}

// RetryableError represents an error that may succeed on retry.
// Examples: network blip, rate limit, server error.
type RetryableError struct {
	Err     error
	Message string
}

func (e *RetryableError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

// Slack API error codes that indicate permanent failures.
var permanentErrorCodes = map[string]bool{
	"invalid_auth":           true,
	"token_revoked":          true,
	"token_expired":          true,
	"not_authed":             true,
	"account_inactive":       true,
	"no_permission":          true,
	"org_login_required":     true,
	"ekm_access_denied":      true,
	"missing_scope":          true,
	"cannot_auth_team":       true,
	"invalid_app_level_token": true,
	"app_not_installed":      true,
}

// IsPermanentError returns true if the error should not be retried.
func IsPermanentError(err error) bool {
	var permanent *PermanentError
	return errors.As(err, &permanent)
}

// IsRetryableError returns true if the error may succeed on retry.
func IsRetryableError(err error) bool {
	var retryable *RetryableError
	return errors.As(err, &retryable)
}

// ClassifySlackError classifies a Slack API error as permanent or retryable.
func ClassifySlackError(slackError string) error {
	if permanentErrorCodes[slackError] {
		return &PermanentError{Message: "slack error: " + slackError}
	}
	// Rate limits and other errors are retryable
	return &RetryableError{Message: "slack error: " + slackError}
}

// ClassifyHTTPError classifies an HTTP response status as permanent or retryable.
func ClassifyHTTPError(statusCode int) error {
	switch {
	case statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden:
		return &PermanentError{Message: "http error: " + http.StatusText(statusCode)}
	case statusCode == http.StatusTooManyRequests:
		return &RetryableError{Message: "rate limited"}
	case statusCode >= 500:
		return &RetryableError{Message: "server error: " + http.StatusText(statusCode)}
	case statusCode >= 400:
		return &PermanentError{Message: "client error: " + http.StatusText(statusCode)}
	default:
		return nil
	}
}

// ClassifyNetworkError classifies a network error as retryable.
func ClassifyNetworkError(err error) error {
	if err == nil {
		return nil
	}
	// Most network errors are retryable
	errStr := err.Error()
	if strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "reset by peer") ||
		strings.Contains(errStr, "broken pipe") ||
		strings.Contains(errStr, "EOF") {
		return &RetryableError{Err: err, Message: "network error"}
	}
	return &RetryableError{Err: err, Message: "network error"}
}
