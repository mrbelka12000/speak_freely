package v1

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

const (
	maxRequestSize = 100 << 20 // ~100mb
)

func (h *Handler) TranscriptCreate(w http.ResponseWriter, r *http.Request) {
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

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("failed to copy file to buffer")
		return
	}

	//themeIDStr := r.FormValue("theme_id")
	//if themeIDStr == "" {
	//	h.writeError(w, errors.New("theme_id is required"), http.StatusBadRequest)
	//	h.log.With("error", err).Error("failed to get theme id")
	//	return
	//}
	//
	//themeID, err := strconv.ParseInt(themeIDStr, 10, 64)
	//if err != nil {
	//	h.writeError(w, err, http.StatusBadRequest)
	//	h.log.With("error", err).Error("failed to get theme id")
	//	return
	//}
	//
	//languageIDStr := r.FormValue("language_id")
	//if languageIDStr == "" {
	//	h.writeError(w, errors.New("language_id is required"), http.StatusBadRequest)
	//	h.log.With("error", err).Error("failed to get language id")
	//	return
	//}
	//
	//languageID, err := strconv.ParseInt(languageIDStr, 10, 64)
	//if err != nil {
	//	h.writeError(w, err, http.StatusBadRequest)
	//	h.log.With("error", err).Error("failed to get language id")
	//	return
	//}

	fileID, missed, err := h.uc.SaveFile(
		context.WithoutCancel(r.Context()),
		&buf,
		handler.Filename,
		handler.Header.Get("Content-Type"),
		handler.Size,
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

	_ = fileID
	w.WriteHeader(http.StatusCreated)
}
