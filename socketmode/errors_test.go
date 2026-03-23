package socketmode

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestClassifySlackError(t *testing.T) {
	permanentErrors := []string{
		"invalid_auth",
		"token_revoked",
		"token_expired",
		"not_authed",
		"account_inactive",
		"no_permission",
		"org_login_required",
		"ekm_access_denied",
		"missing_scope",
		"cannot_auth_team",
		"invalid_app_level_token",
		"app_not_installed",
	}

	for _, code := range permanentErrors {
		t.Run(code+" is permanent", func(t *testing.T) {
			err := ClassifySlackError(code)
			if !IsPermanentError(err) {
				t.Errorf("expected %s to be permanent error", code)
			}
		})
	}

	retryableErrors := []string{
		"rate_limited",
		"service_unavailable",
		"internal_error",
		"unknown_error",
	}

	for _, code := range retryableErrors {
		t.Run(code+" is retryable", func(t *testing.T) {
			err := ClassifySlackError(code)
			if !IsRetryableError(err) {
				t.Errorf("expected %s to be retryable error", code)
			}
		})
	}
}

func TestClassifyHTTPError(t *testing.T) {
	tests := []struct {
		status    int
		permanent bool
		retryable bool
	}{
		{http.StatusUnauthorized, true, false},
		{http.StatusForbidden, true, false},
		{http.StatusNotFound, true, false},
		{http.StatusBadRequest, true, false},
		{http.StatusTooManyRequests, false, true},
		{http.StatusInternalServerError, false, true},
		{http.StatusBadGateway, false, true},
		{http.StatusServiceUnavailable, false, true},
		{http.StatusGatewayTimeout, false, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("status %d", tt.status), func(t *testing.T) {
			err := ClassifyHTTPError(tt.status)
			if tt.permanent && !IsPermanentError(err) {
				t.Errorf("expected status %d to be permanent error", tt.status)
			}
			if tt.retryable && !IsRetryableError(err) {
				t.Errorf("expected status %d to be retryable error", tt.status)
			}
		})
	}

	t.Run("status 200 returns nil", func(t *testing.T) {
		err := ClassifyHTTPError(http.StatusOK)
		if err != nil {
			t.Errorf("expected nil for status 200, got %v", err)
		}
	})
}

func TestClassifyNetworkError(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		err := ClassifyNetworkError(nil)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	networkErrors := []string{
		"connection refused",
		"no such host",
		"timeout",
		"reset by peer",
		"broken pipe",
		"EOF",
	}

	for _, msg := range networkErrors {
		t.Run(msg+" is retryable", func(t *testing.T) {
			err := ClassifyNetworkError(errors.New(msg))
			if !IsRetryableError(err) {
				t.Errorf("expected %q to be retryable error", msg)
			}
		})
	}

	t.Run("wraps original error", func(t *testing.T) {
		original := errors.New("connection refused")
		classified := ClassifyNetworkError(original)

		var retryable *RetryableError
		if !errors.As(classified, &retryable) {
			t.Fatal("expected RetryableError")
		}

		if !errors.Is(retryable.Err, original) {
			t.Error("expected original error to be wrapped")
		}
	})
}

func TestErrorHelpers(t *testing.T) {
	t.Run("IsPermanentError with PermanentError", func(t *testing.T) {
		err := &PermanentError{Message: "test"}
		if !IsPermanentError(err) {
			t.Error("expected true for PermanentError")
		}
	})

	t.Run("IsPermanentError with other error", func(t *testing.T) {
		err := errors.New("random error")
		if IsPermanentError(err) {
			t.Error("expected false for non-PermanentError")
		}
	})

	t.Run("IsPermanentError with wrapped PermanentError", func(t *testing.T) {
		inner := &PermanentError{Message: "inner"}
		wrapped := fmt.Errorf("outer: %w", inner)
		if !IsPermanentError(wrapped) {
			t.Error("expected true for wrapped PermanentError")
		}
	})

	t.Run("IsRetryableError with RetryableError", func(t *testing.T) {
		err := &RetryableError{Message: "test"}
		if !IsRetryableError(err) {
			t.Error("expected true for RetryableError")
		}
	})

	t.Run("IsRetryableError with other error", func(t *testing.T) {
		err := errors.New("random error")
		if IsRetryableError(err) {
			t.Error("expected false for non-RetryableError")
		}
	})

	t.Run("IsRetryableError with wrapped RetryableError", func(t *testing.T) {
		inner := &RetryableError{Message: "inner"}
		wrapped := fmt.Errorf("outer: %w", inner)
		if !IsRetryableError(wrapped) {
			t.Error("expected true for wrapped RetryableError")
		}
	})
}

func TestPermanentError(t *testing.T) {
	t.Run("Error with message only", func(t *testing.T) {
		err := &PermanentError{Message: "invalid token"}
		if err.Error() != "invalid token" {
			t.Errorf("got %q, want %q", err.Error(), "invalid token")
		}
	})

	t.Run("Error with wrapped error", func(t *testing.T) {
		inner := errors.New("inner error")
		err := &PermanentError{Err: inner, Message: "outer"}
		if err.Error() != "outer: inner error" {
			t.Errorf("got %q, want %q", err.Error(), "outer: inner error")
		}
	})

	t.Run("Unwrap returns inner error", func(t *testing.T) {
		inner := errors.New("inner error")
		err := &PermanentError{Err: inner, Message: "outer"}
		if !errors.Is(err, inner) {
			t.Error("expected Unwrap to return inner error")
		}
	})
}

func TestRetryableError(t *testing.T) {
	t.Run("Error with message only", func(t *testing.T) {
		err := &RetryableError{Message: "network error"}
		if err.Error() != "network error" {
			t.Errorf("got %q, want %q", err.Error(), "network error")
		}
	})

	t.Run("Error with wrapped error", func(t *testing.T) {
		inner := errors.New("connection refused")
		err := &RetryableError{Err: inner, Message: "network error"}
		if err.Error() != "network error: connection refused" {
			t.Errorf("got %q, want %q", err.Error(), "network error: connection refused")
		}
	})

	t.Run("Unwrap returns inner error", func(t *testing.T) {
		inner := errors.New("connection refused")
		err := &RetryableError{Err: inner, Message: "network error"}
		if !errors.Is(err, inner) {
			t.Error("expected Unwrap to return inner error")
		}
	})
}

func TestSentinelErrors(t *testing.T) {
	sentinels := []error{
		ErrHelloTimeout,
		ErrConnectionClosed,
		ErrWriteTimeout,
		ErrShuttingDown,
		ErrHandlerTimeout,
		ErrConcurrencyLimit,
	}

	for _, err := range sentinels {
		t.Run(err.Error(), func(t *testing.T) {
			if err == nil {
				t.Error("sentinel error should not be nil")
			}
			if err.Error() == "" {
				t.Error("sentinel error should have message")
			}
		})
	}
}
