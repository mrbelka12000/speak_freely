package service

import (
	"github.com/mrbelka12000/linguo_sphere_backend/internal/repository"
	userservice "github.com/mrbelka12000/linguo_sphere_backend/internal/service/user"
)

type Service struct {
	User *userservice.Service
}

func New(r *repository.Repo) *Service {
	return &Service{
		User: userservice.New(r.User),
	}
}
