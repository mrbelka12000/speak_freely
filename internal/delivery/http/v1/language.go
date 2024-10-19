package v1

import (
	"encoding/json"
	"net/http"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

func (h *Handler) LanguageCreate(w http.ResponseWriter, r *http.Request) {
	var obj models.LanguageCU
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Info("can not decode")
		return
	}

	missed, err := h.uc.LanguageValidate(r.Context(), obj)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		h.log.With("error", err).Info("can not validate language")
		return
	}
	if len(missed) > 0 {
		writeJson(w, missed, http.StatusBadRequest)
		return
	}

	err = h.uc.LanguageCreate(r.Context(), obj)
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
	langs, count, err := h.uc.LanguageList(r.Context())
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		h.log.With("error", err).Info("can not list languages")
		return
	}

	writeJson(w, models.PaginatedResponse{
		Result: langs,
		Page:   1,
		Count:  count,
	}, http.StatusOK)
}
