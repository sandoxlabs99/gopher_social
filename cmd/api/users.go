package main

import (
	"context"
	"errors"
	"fmt"
	"gopher_social/internal/models"
	"gopher_social/internal/store"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ctxKeyUser string

const UserContextKey = ctxKeyUser("user")

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	models.User
//	@Failure		400		{object}	error	"Invalid user ID"
//	@Failure		401		{object}	error	"Unauthorized"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.JSONResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{object}	string	"Follow successful"
//	@Failure		400		{object}	error	"Invalid payload"
//	@Failure		401		{object}	error	"Unauthorized"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		409		{object}	error	"Already following user"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser, err := getAuthUserFromContext(r) // the authenticated user
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	followedUser := getUserFromCtx(r) // from user params

	if followerUser.ID == followedUser.ID {
		app.badRequestResponse(w, r, fmt.Errorf("cannot follow yourself"))
		return
	}

	err = app.store.Followers.Follow(r.Context(), followerUser.ID, followedUser.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateKey):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnfollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{object}	string	"Unfollow successful"
//	@Failure		400		{object}	error	"Invalid payload"
//	@Failure		401		{object}	error	"Unauthorized"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowerUser, err := getAuthUserFromContext(r) // the authenticated user
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	unfollowedUser := getUserFromCtx(r) // from user params

	if unfollowerUser.ID == unfollowedUser.ID {
		app.badRequestResponse(w, r, fmt.Errorf("cannot unfollow yourself"))
		return
	}

	err = app.store.Followers.UnFollow(r.Context(), unfollowedUser.ID, unfollowedUser.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ActivateUser godoc
//
//	@Summary		Activate a user account
//	@Description	Activate a user account using the provided token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		400		{object}	error	"Invalid token"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if err := app.store.Users.Activate(r.Context(), token); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err, "user not found")
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// middleware
// get users middleware
// gets user from params and check if the user exists and returns a user
func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.getUser(ctx, userID)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err, "user not found")
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, UserContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *models.User {
	user, _ := r.Context().Value(UserContextKey).(*models.User)

	return user
}
