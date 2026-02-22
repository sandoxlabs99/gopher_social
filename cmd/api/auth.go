package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gopher_social/internal/mailer"
	"gopher_social/internal/models"
	"gopher_social/internal/store"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required,max=15"`
	LastName  string `json:"lastName" validate:"required,max=15"`
	Username  string `json:"username" validate:"required,min=3,max=100"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=8,max=72"`
}

type UserWithToken struct {
	*models.User
	Token string `json:"token"`
}

// registerUser godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with the provided details
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User registration details"
//	@Success		201		{object}	UserWithToken
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &models.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Username:  payload.Username,
		Email:     payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	token, err := uuid.NewV7()
	if err != nil {
		app.internalServerError(w, r, err)
	}

	hash := sha256.Sum256([]byte(token.String()))
	hashToken := hex.EncodeToString(hash[:])

	ctx := r.Context()

	// store the user
	if err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	UserWithToken := UserWithToken{
		User:  user,
		Token: token.String(),
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, token.String())

	isProdEnv := app.config.namespace == "production"

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      UserWithToken.Username,
		ActivationURL: activationURL,
	}

	// send the email invite
	err = app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw(
			"error sending welcome email",
			"error", err,
		)

		// rollback user creation if email fails (SAGA Pattern)
		if err := app.store.Users.Delete(ctx, UserWithToken.ID); err != nil {
			app.logger.Errorw(
				"error deleting user",
				"error", err,
			)
		}

		app.internalServerError(w, r, err)
		return
	}

	if err := app.JSONResponse(w, http.StatusCreated, UserWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// createTokenHandler godoc
//
//	@Summary		Creates an auth token
//	@Description	Generates a new auth token for the user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials for token generation"
//	@Success		201		{object}	string					"Authentication token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// parse the payload credentials
	var payload CreateUserTokenPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// fetch the user (check if user exists) from the payload
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unAuthorizedErrorResponse(w, r, fmt.Errorf("invalid email or password"))
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	user.Password.Text = &payload.Password

	if err := user.Password.Verify(); err != nil {
		app.unAuthorizedErrorResponse(w, r, fmt.Errorf("invalid email or password"))
		return
	}

	// generate the token -> add claims
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(user.ID, 10),
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.JSONResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
