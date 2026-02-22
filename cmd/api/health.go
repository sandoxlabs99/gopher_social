package main

import (
	"net/http"
)

// CheckHealth godoc
//
//	@Summary		Tests API's health endpoint
//	@Description	Checks the API's status by sending some info about it.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	main.JSONResponse.envelope
//	@Failure		500	{object}	main.writeJSONError.envelope
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]any{
		"version": version,
		"status":  "ok",
		"env":     app.config.namespace,
	}

	// if err := errors.New("a test error"); err != nil {
	// 	writeJSONError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	if err := app.JSONResponse(w, http.StatusOK, result); err != nil {
		app.internalServerError(w, r, err)
	}
	// w.Write([]byte("Health check: OK"))
}
