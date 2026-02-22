package main

import (
	"gopher_social/internal/auth"
	"gopher_social/internal/ratelimiter"
	"gopher_social/internal/store"
	"gopher_social/internal/store/cache"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	// logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockRedisStore := cache.NewMockRedisStorage()
	testAuth := &auth.TestAuthenticator{}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	return &application{
		config:        cfg,
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockRedisStore,
		authenticator: testAuth,
		rateLimiter:   rateLimiter,
	}
}

func executeRequest(mux http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
