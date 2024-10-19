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
	key, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		panic("no secret key provided")
	}
	secretKey = []byte(key)

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
	r.Use(h.recovery)
	r.Use(h.cors)
	// users
	r.HandleFunc("/api/v1/register", h.Registration).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/login", h.Login).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/confirm", h.ConfirmEmail)
	r.HandleFunc("/api/v1/profile", h.authenticateMiddleware(h.UpdateProfile, true)).Methods(http.MethodPut)
	r.HandleFunc("/api/v1/profile", h.authenticateMiddleware(h.Profile, true)).Methods(http.MethodGet)

	// languages
	r.HandleFunc("/api/v1/lang", h.authenticateMiddleware(h.LanguageCreate, true)).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/langs", h.LanguageList)

	//themes
	r.HandleFunc("/api/v1/theme/{id}", h.authenticateMiddleware(h.GetTheme, false)).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/theme", h.authenticateMiddleware(h.CreateTheme, true)).Methods(http.MethodPost)

	// transcripts
	r.HandleFunc("/api/v1/transcript", h.TranscriptCreate)
	// tokens
	r.HandleFunc("/api/v1/tokens", h.Tokens)
}

func (h *Handler) writeError(w http.ResponseWriter, err error, code int) {
	errResp := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	writeJson(w, errResp, code)
}

func writeJson(w http.ResponseWriter, data interface{}, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	body, _ := json.Marshal(data)
	w.Write(body)
}
