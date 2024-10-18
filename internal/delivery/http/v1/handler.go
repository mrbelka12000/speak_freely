package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/usecase"
)

type (
	Handler struct {
		uc  *usecase.UseCase
		log *slog.Logger
	}
)

// New
func New(uc *usecase.UseCase, opts ...opt) *Handler {
	h := &Handler{
		uc:  uc,
		log: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// InitRoutes
func (h *Handler) InitRoutes(r *mux.Router) {
	r.HandleFunc("/register", h.Registration)
	r.HandleFunc("/login", h.Login)
	r.HandleFunc("/confirm", h.ConfirmEmail)
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
