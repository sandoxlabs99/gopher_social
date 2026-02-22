package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw(
		"internal server error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	writeJSONError(w, http.StatusInternalServerError, "something went wrong")
}
func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw(
		"forbidden error",
		"method", r.Method,
		"path", r.URL.Path,
	)
	writeJSONError(w, http.StatusForbidden, "forbidden to access this resource")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw(
		"bad request error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error, msg string) {
	app.logger.Warnw(
		"not found error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	writeJSONError(w, http.StatusNotFound, msg)
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw(
		"conflict error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unAuthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw(
		"unauthorized basic error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (app *application) unAuthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw(
		"unauthorized error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after:"+retryAfter)
}
