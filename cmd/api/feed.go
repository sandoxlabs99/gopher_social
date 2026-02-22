package main

import (
	"gopher_social/internal/utils"
	"net/http"
)

// GetUserFeed godoc
//
//	@Summary		Fetches the authenticated users feed
//	@Description	Fetches the authenticated users feed with optional filters
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int			false	"Number of posts to return"		default(20)		minimum(1)	maximum(100)
//	@Param			offset	query		int			false	"Number of posts to skip"		default(0)		minimum(0)
//	@Param			sort	query		string		false	"Sort order: 'asc' or 'desc'"	default(desc)	enum(asc,desc)
//	@Param			search	query		string		false	"Search term to filter posts by title or content"
//	@Param			tags	query		[]string	false	"Filter posts by tags"	minItems(1)	maxItems(5)	items(minLength(2),maxLength(30))
//	@Success		200		{object}	[]models.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := utils.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Search: "",
		Tags:   []string{},
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	userID := int64(25)
	feed, err := app.store.Posts.GetUserFeed(r.Context(), userID, fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.JSONResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
