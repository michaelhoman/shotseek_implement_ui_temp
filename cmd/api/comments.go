package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	store "github.com/michaelhoman/ShotSeek/internal/store/postgres"
)

// type commentKey string

// const commentCtx commentKey = "comment"

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := store.Comment{
		PostID:  postID,
		UserID:  1,
		Content: payload.Content,
	}

	ctx := r.Context()

	commentsStore := app.store.Comments
	if err := commentsStore.Create(ctx, &comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// deleteCommentByPostHandler deletes comments by post ID
func (app *application) DeleteByPostID(w http.ResponseWriter, r *http.Request) {

	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	post, err := app.store.Posts.GetByID(r.Context(), postID)
	if err != nil {
		app.notFoundResponse(w, r, err)
		return
	}

	ctx := r.Context()

	commentsStore := app.store.Comments
	if err := commentsStore.DeleteByPostID(ctx, post.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
