package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

const (
	maxRequestSize = 100 << 20 // ~100mb
)

func (h *Handler) TranscriptCreate(w http.ResponseWriter, r *http.Request) {
	uAny := r.Context().Value(userObj)
	user, ok := uAny.(models.User)
	if !ok {
		h.writeError(w, errors.New("invalid user"), http.StatusBadRequest)
		h.log.Error("invalid user")
		return
	}

	err := r.ParseMultipartForm(maxRequestSize)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("failed to parse multipart form")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("failed to get file from form")
		return
	}
	defer file.Close()

	if handler.Header.Get("Content-Type") == "" {
		h.writeError(w, fmt.Errorf("content type is required"), http.StatusBadRequest)
		h.log.With("error", err).Error("content type is required")
		return
	}

	f, err := handler.Open()
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("failed to get file from handler")
		return
	}

	themeIDStr := r.FormValue("theme_id")
	if themeIDStr == "" {
		h.writeError(w, errors.New("theme_id is required"), http.StatusBadRequest)
		h.log.With("error", err).Error("failed to get theme id")
		return
	}

	themeID, err := strconv.ParseInt(themeIDStr, 10, 64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("failed to get theme id")
		return
	}

	languageIDStr := r.FormValue("language_id")
	if languageIDStr == "" {
		h.writeError(w, errors.New("language_id is required"), http.StatusBadRequest)
		h.log.With("error", err).Error("failed to get language id")
		return
	}

	languageID, err := strconv.ParseInt(languageIDStr, 10, 64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("failed to get language id")
		return
	}

	id, missed, err := h.uc.TranscriptBuild(
		context.WithoutCancel(r.Context()),
		f,
		handler.Filename,
		handler.Header.Get("Content-Type"),
		handler.Size,
		languageID,
		themeID,
		user.ID,
	)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		h.log.With("error", err).Error("failed to upload file")
		return
	}
	if len(missed) > 0 {
		writeJson(w, missed, http.StatusBadRequest)
		return
	}

	writeJson(w, struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	}, http.StatusCreated)
}

func (h *Handler) TranscriptGet(w http.ResponseWriter, r *http.Request) {
	uAny := r.Context().Value(userObj)
	user, ok := uAny.(models.User)
	if !ok {
		h.writeError(w, errors.New("invalid users"), http.StatusBadRequest)
		h.log.Error("invalid user")
		return
	}

	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("failed to get id")
		return
	}

	obj, err := h.uc.TranscriptGet(r.Context(), id, user)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("failed to get transcript")
		return
	}

	writeJson(w, obj, http.StatusOK)
}

func (h *Handler) TranscriptList(w http.ResponseWriter, r *http.Request) {
	var pars models.TranscriptListPars
	err := h.decoder.Decode(&pars, r.URL.Query())
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("can not decode pars")
		return
	}

	pars.PaginationParams, err = getPaginationParams(pars.PaginationParams)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("can not get pagination params")
		return
	}

	result, count, err := h.uc.TranscriptList(r.Context(), pars)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		h.log.With("error", err).Error("failed to list transcripts")
		return
	}

	writeJson(w, models.PaginatedResponse{
		Result: result,
		Page:   pars.PaginationParams.Page,
		Count:  count,
	}, http.StatusOK)
}
