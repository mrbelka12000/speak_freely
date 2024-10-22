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

	id, missed, err := h.uc.UserCreate(r.Context(), obj)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		h.log.With("error", err).Error("can not create user")
		return
	}

	if len(missed) > 0 {
		writeJson(w, missed, http.StatusBadRequest)
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
	uAny := r.Context().Value(userObj)
	user, ok := uAny.(models.User)
	if !ok {
		h.writeError(w, errors.New("invalid user"), http.StatusBadRequest)
		h.log.Error("invalid user")
		return
	}

	var obj models.UserCU
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("can not decode")
		return
	}

	missed, err := h.uc.UserUpdate(r.Context(), models.UserGet{
		ID: user.ID,
	}, obj)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("can not update user")
		return
	}

	if len(missed) > 0 {
		writeJson(w, missed, http.StatusBadRequest)
		return
	}

	writeJson(w, struct {
		Success bool `json:"success"`
	}{
		Success: true,
	}, http.StatusOK)
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	uAny := r.Context().Value(userObj)
	user, ok := uAny.(models.User)
	if !ok {
		h.writeError(w, errors.New("invalid user"), http.StatusBadRequest)
		h.log.Error("invalid user")
		return
	}

	user, err := h.uc.UserGet(r.Context(), models.UserGet{ID: user.ID})
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		h.log.With("error", err).Error("can not get user")
		return
	}

	writeJson(w, user, http.StatusOK)
}

func (h *Handler) UsersList(w http.ResponseWriter, r *http.Request) {

	var pars models.UserListPars
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

	result, count, err := h.uc.UserList(r.Context(), pars)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		h.log.With("error", err).Error("can not get users")
		return
	}

	writeJson(w, models.PaginatedResponse{
		Result: result,
		Page:   pars.PaginationParams.Page,
		Count:  count,
	}, http.StatusOK)
}
