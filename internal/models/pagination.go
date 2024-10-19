package models

type (
	PaginationParams struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Page   int `json:"page"`
	}

	PaginatedResponse struct {
		Result any `json:"result"`
		Page   int `json:"page"`
		Count  int `json:"count"`
	}
)
