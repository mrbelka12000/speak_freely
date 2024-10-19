package v1

import (
	"context"
	"fmt"
	"net/http"
)

type (
	ctxKey string
)

const (
	id ctxKey = "id"
)

func (h *Handler) recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				h.writeError(w, fmt.Errorf("%v", err), http.StatusInternalServerError)
				h.log.Error(fmt.Sprintf("panic: %v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) authenticateMiddleware(next http.HandlerFunc, strict bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" && strict {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if strict {
			token, err := verifyToken(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*myClaims)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				h.log.Error("can not convert claims")
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), id, claims.ID))
		}

		next.ServeHTTP(w, r)
	}
}
