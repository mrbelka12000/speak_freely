package v1

import (
	"encoding/json"
	"net/http"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

func (h *Handler) LanguageCreate(w http.ResponseWriter, r *http.Request) {
	var obj models.Language
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Info("can not decode")
		return
	}

	err = h.uc.LanguageCreate(r.Context(), obj.Name)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		h.log.With("error", err).Info("can not create language")
		return
	}

	writeJson(w, struct {
		Success bool `json:"success"`
	}{
		Success: true,
	}, http.StatusOK)
}

func (h *Handler) LanguageList(w http.ResponseWriter, r *http.Request) {
	langs, _, err := h.uc.LanguageList(r.Context())
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		h.log.With("error", err).Info("can not list languages")
		return
	}

	writeJson(w, langs, http.StatusOK)
}
