package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/usecase"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/validate"
)

type (
	Handler struct {
		uc  *usecase.UseCase
		v   *validate.Validator
		log *slog.Logger
	}
)

func New(uc *usecase.UseCase, v *validate.Validator) Handler {
	return Handler{
		uc: uc,
		v:  v,
	}
}

func (h *Handler) InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/register", h.CreateUser)
}

func (h *Handler) writeBadRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
	h.log.Error(err.Error())
}

func (h *Handler) writeInternalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	h.log.Error(err.Error())
}

func writeError(w http.ResponseWriter, err error, code int) {
	errResp := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	writeJson(w, errResp, code)
}

func writeJson(w http.ResponseWriter, data interface{}, httpStatus int) {
	w.WriteHeader(httpStatus)
	w.Header().Set("Content-Type", "application/json")
	body, _ := json.Marshal(data)
	w.Write(body)
}
