package v1

import (
	"encoding/json"
	"net/http"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var obj models.UserCU

	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		h.writeBadRequest(w, err)
		return
	}

	missed, err := h.v.ValidateUser(r.Context(), obj, -1)
	if err != nil {
		h.writeInternalServerError(w, err)
		return
	}
	if len(missed) > 0 {
		writeJson(w, missed, http.StatusBadRequest)
		return
	}

	id, err := h.uc.UserCreate(r.Context(), obj)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJson(w, struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	}, http.StatusCreated)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var obj models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		h.writeBadRequest(w, err)
		return
	}

	id, err := h.uc.UserLogin(r.Context(), obj)
	if err != nil {
		h.writeBadRequest(w, err)
		return
	}

	writeJson(w, struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	}, http.StatusOK)
}
