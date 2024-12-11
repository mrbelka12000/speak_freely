package service

import (
	"github.com/mrbelka12000/speak_freely/internal/repository"
	fileservice "github.com/mrbelka12000/speak_freely/internal/service/file"
	languageservice "github.com/mrbelka12000/speak_freely/internal/service/language"
	themeservice "github.com/mrbelka12000/speak_freely/internal/service/theme"
	topicservice "github.com/mrbelka12000/speak_freely/internal/service/topic"
	transcriptservice "github.com/mrbelka12000/speak_freely/internal/service/transcript"
	userservice "github.com/mrbelka12000/speak_freely/internal/service/user"
)

// Service adapter for usecase
type Service struct {
	User       *userservice.Service
	Language   *languageservice.Service
	Theme      *themeservice.Service
	File       *fileservice.Service
	Transcript *transcriptservice.Service
	Topic      *topicservice.Service
}

// New create instance of service
func New(r *repository.Repo) *Service {
	return &Service{
		User:       userservice.New(r.User),
		Language:   languageservice.New(r.Language),
		Theme:      themeservice.New(r.Theme),
		File:       fileservice.New(r.File),
		Transcript: transcriptservice.New(r.Transcript),
		Topic:      topicservice.New(r.Topic),
	}
}
