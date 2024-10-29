package v1

import (
	"errors"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

func getPaginationParams(pp models.PaginationParams) (models.PaginationParams, error) {
	if pp.Page <= 0 {
		pp.Page = 1
	}

	if pp.Limit <= 0 {
		return pp, errors.New("limit cannot be zero or negative")
	}

	return models.PaginationParams{
		Limit:  pp.Limit,
		Offset: (pp.Page - 1) * pp.Limit,
		Page:   pp.Page,
	}, nil
}
