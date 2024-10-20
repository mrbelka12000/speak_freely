package models

type (
	PaginationParams struct {
		Limit  int `json:"limit" schema:"limit"`
		Offset int `json:"offset" schema:"offset"`
		Page   int `json:"page" schema:"page"`
	}

	PaginatedResponse struct {
		Result any `json:"result"`
		Page   int `json:"page"`
		Count  int `json:"count"`
	}
)
