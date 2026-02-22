package main

import (
	"fmt"
	"gopher_social/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CreateCommentsPayload struct {
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required,min=3,max=1000"`
}

// CreateComment godoc
//
//	@Summary		Creates a comment on a post
//	@Description	Creates a comment on a post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int64					true	"Post ID"
//	@Param			payload	body		CreateCommentsPayload	true	"Comment payload"
//	@Success		201		{object}	models.Comment
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/comments [post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")

	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	var payload CreateCommentsPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid payload"))
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &models.Comment{
		PostID:  postID,
		UserID:  payload.UserID,
		Content: payload.Content,
	}

	if err := app.store.Comments.Create(r.Context(), comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.JSONResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
