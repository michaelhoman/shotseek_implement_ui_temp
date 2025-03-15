package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	// store "github.com/michaelhoman/ShotSeek/internal/store/postgres"
	"github.com/michaelhoman/ShotSeek/internal/store"
	"github.com/michaelhoman/ShotSeek/internal/utils"
)

type commentKey string

const commentCtx commentKey = "comment"

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

// CreateComment godoc
//
//	@Summary		Creates a comment
//	@Description	Creates a new comment from payload
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Post ID"
//	@Param			payload	body		CreateCommentPayload	true	"Comment payload"
//	@Success		200		{object}	store.Comment
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id}/comments [post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
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
		utils.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// GetComment godoc
//
//	@Summary		Retrieves a comment
//	@Description	Retrieves a comment by ID
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Comment ID"
//	@Success		200	{object}	store.Comment
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/comments/{id} [get]
func (app *application) getCommentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getCommentHandler")
	comment := getCommentFromCtx(r)
	fmt.Println("comment:", comment)

	if err := utils.JsonResponse(w, http.StatusOK, comment); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

// UpdateComment godoc
//
//	@Summary		Updates a comment
//	@Description	Updates a comment by ID
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Comment ID"
//	@Param			payload	body		CreateCommentPayload	true	"Comment payload"
//	@Success		200		{object}	store.Comment
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/comments/{id} [patch]
func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	comment := getCommentFromCtx(r)

	var payload CreateCommentPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	comment.Content = payload.Content

	ctx := r.Context()

	commentsStore := app.store.Comments
	if err := commentsStore.Update(ctx, comment); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := utils.JsonResponse(w, http.StatusOK, comment); err != nil {
		utils.InternalServerError(w, r, err)
	}
}

// DeleteComment godoc
//
//	@Summary		Deletes a comment
//	@Description	Deletes a comment by ID
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Comment ID"
//	@Success		204	{object}	nil
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/comments/{id} [delete]
func (app *application) DeleteByCommentIDHandler(w http.ResponseWriter, r *http.Request) {

	commentID, err := strconv.ParseInt(chi.URLParam(r, "commentID"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	comment, err := app.store.Comments.GetByCommentID(r.Context(), commentID)
	if err != nil {
		utils.NotFoundResponse(w, r, err)
		return
	}

	ctx := r.Context()

	commentsStore := app.store.Comments
	if err := commentsStore.DeleteByCommentID(ctx, comment.ID); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) DeleteByPostID(w http.ResponseWriter, r *http.Request) {

	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	post, err := app.store.Posts.GetByID(r.Context(), postID)
	if err != nil {
		utils.NotFoundResponse(w, r, err)
		return
	}

	ctx := r.Context()

	commentsStore := app.store.Comments
	if err := commentsStore.DeleteByPostID(ctx, post.ID); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) commentsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("commentsContextMiddleware")
		idParam := chi.URLParam(r, "commentID")
		id, err := strconv.ParseInt(idParam, 10, 64)

		if err != nil {
			fmt.Println("Error in commentsContextMiddleware")
			utils.InternalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		comment, err := app.store.Comments.GetByCommentID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				utils.NotFoundResponse(w, r, err)
			default:
				utils.InternalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, commentCtx, comment)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getCommentFromCtx(r *http.Request) *store.Comment {
	comment, _ := r.Context().Value(commentCtx).(*store.Comment)
	return comment
}
