package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	store "github.com/michaelhoman/ShotSeek/internal/store/postgres"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostsHandler(w http.ResponseWriter, r *http.Request) {

	// Step 1: Read the body to log it
	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	app.badRequestResponse(w, r, err)
	// 	return
	// }
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// fmt.Println("Request body:", string(body))

	// Step 2: Reset the body so readJSON can decode it
	// r.Body = io.NopCloser(bytes.NewReader(body))

	// Step 3: Decode the JSON into the payload
	// if err := readJSON(w, r, &payload); err != nil {
	// 	app.badRequestResponse(w, r, err)
	// 	return
	// }

	post := store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	fmt.Println("Tags data:", post.Tags)

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, &post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, postID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)

		}
		return
	}

	comments, err := app.store.Comments.GetByPostID(ctx, postID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// if post == nil {
	// 	writeJSONError(w, http.StatusNotFound, err.Error())
	// 	return
	// }

	// fmt.Println("Tags data:", post.Tags)
	// fmt.Println("Post data:", post)

	// if err != nil {
	// 	writeJSONError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// if err := writeJSON(w, http.StatusOK, post); err != nil {
	// 	writeJSONError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
	ctx := r.Context()

	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	post := store.Post{
		ID:      postID,
		Content: payload.Content,
		Title:   payload.Title,
		Tags:    payload.Tags,
	}

	fmt.Println("Post data:", post)
	app.store.Posts.UpdateByID(ctx, post)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Posts.UpdateByID(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	ctx := r.Context()
	if err := app.store.Comments.DeleteByPostID(ctx, postID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.store.Posts.DeleteByID(ctx, postID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
