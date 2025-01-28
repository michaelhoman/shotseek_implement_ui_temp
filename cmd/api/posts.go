package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	postgres_store "github.com/michaelhoman/ShotSeek/internal/store/postgres"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required, max=100"`
	Content string   `json:"content" validate:"required, max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostsHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	// Step 1: Read the body to log it
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("Request body:", string(body))

	// Step 2: Reset the body so readJSON can decode it
	r.Body = io.NopCloser(bytes.NewReader(body))

	// Step 3: Decode the JSON into the payload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &postgres_store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	fmt.Println("Tags data:", post.Tags)

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
