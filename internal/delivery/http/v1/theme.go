package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

func (h *Handler) CreateTheme(w http.ResponseWriter, r *http.Request) {
	var obj models.ThemeCU
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("decode request body error")
		return
	}

	id, missed, err := h.uc.ThemeBuild(r.Context(), obj)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		h.log.With("error", err).Error("build theme")
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

func (h *Handler) GetTheme(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("parse theme id")
		return
	}

	obj, err := h.uc.ThemeGet(r.Context(), id)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("get theme")
		return
	}

	writeJson(w, obj, http.StatusOK)
}
