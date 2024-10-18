package service

import (
	"github.com/mrbelka12000/linguo_sphere_backend/internal/repository"
	languageservice "github.com/mrbelka12000/linguo_sphere_backend/internal/service/language"
	userservice "github.com/mrbelka12000/linguo_sphere_backend/internal/service/user"
)

// Service adapter for usecase
type Service struct {
	User     *userservice.Service
	Language *languageservice.Service
}

// New create instance of service
func New(r *repository.Repo) *Service {
	return &Service{
		User:     userservice.New(r.User),
		Language: languageservice.New(r.Language),
	}
}
