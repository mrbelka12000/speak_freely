package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var obj models.UserCU

	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		h.writeBadRequest(w, err)
		return
	}

	missed, err := h.uc.UserCUValidate(r.Context(), obj, -1)
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		h.writeBadRequest(w, errors.New("code is required"))
		return
	}

	err := h.uc.UserConfirm(r.Context(), code)
	if err != nil {
		h.writeBadRequest(w, err)
		return
	}

	writeJson(w, struct {
		Success bool `json:"success"`
	}{
		Success: true,
	}, http.StatusOK)
}
