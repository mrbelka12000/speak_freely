package service

import (
	"github.com/mrbelka12000/linguo_sphere_backend/internal/repository"
	fileservice "github.com/mrbelka12000/linguo_sphere_backend/internal/service/file"
	languageservice "github.com/mrbelka12000/linguo_sphere_backend/internal/service/language"
	themeservice "github.com/mrbelka12000/linguo_sphere_backend/internal/service/theme"
	transcriptservice "github.com/mrbelka12000/linguo_sphere_backend/internal/service/transcript"
	userservice "github.com/mrbelka12000/linguo_sphere_backend/internal/service/user"
)

// Service adapter for usecase
type Service struct {
	User       *userservice.Service
	Language   *languageservice.Service
	Theme      *themeservice.Service
	File       *fileservice.Service
	Transcript *transcriptservice.Service
}

// New create instance of service
func New(r *repository.Repo) *Service {
	return &Service{
		User:       userservice.New(r.User),
		Language:   languageservice.New(r.Language),
		Theme:      themeservice.New(r.Theme),
		File:       fileservice.New(r.File),
		Transcript: transcriptservice.New(r.Transcript),
	}
}
