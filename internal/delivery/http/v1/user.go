package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

// Registration
func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var obj models.UserCU

	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("can not decode")
		return
	}

	missed, err := h.uc.UserCUValidate(r.Context(), obj, -1)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		h.log.With("error", err).Error("can not validate")
		return
	}
	if len(missed) > 0 {
		writeJson(w, missed, http.StatusBadRequest)
		return
	}

	id, err := h.uc.UserCreate(r.Context(), obj)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		h.log.With("error", err).Error("can not create user")
		return
	}

	writeJson(w, struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	}, http.StatusCreated)
}

// Login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var obj models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("can not decode")
		return
	}

	id, err := h.uc.UserLogin(r.Context(), obj)
	if err != nil {
		h.writeError(w, err, http.StatusUnauthorized)
		h.log.With("error", err).Error("can not login")
		return
	}

	tokens, err := createTokens(id)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		h.log.With("error", err).Error("can not create tokens")
		return
	}

	writeJson(w, tokens, http.StatusOK)
}

// ConfirmEmail
func (h *Handler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		h.writeError(w, errors.New("code is required"), http.StatusBadRequest)
		h.log.Error("code is required")
		return
	}

	err := h.uc.UserConfirm(r.Context(), code)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("can not confirm")
		return
	}

	writeJson(w, struct {
		Success bool `json:"success"`
	}{
		Success: true,
	}, http.StatusOK)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	idAny := r.Context().Value(id)
	idInt64, ok := idAny.(int64)
	if !ok {
		h.writeError(w, errors.New("invalid id"), http.StatusBadRequest)
		h.log.Error("invalid id")
		return
	}

	var obj models.UserCU
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("can not decode")
		return
	}

	err = h.uc.UserUpdate(r.Context(), idInt64, obj)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("can not update user")
		return
	}

	writeJson(w, struct {
		Success bool `json:"success"`
	}{
		Success: true,
	}, http.StatusOK)
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	idAny := r.Context().Value(id)
	idInt64, ok := idAny.(int64)
	if !ok {
		h.writeError(w, errors.New("invalid id"), http.StatusBadRequest)
		h.log.Error("invalid id")
		return
	}

	user, err := h.uc.UserGet(r.Context(), idInt64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("can not get user")
		return
	}

	writeJson(w, user, http.StatusOK)
}
