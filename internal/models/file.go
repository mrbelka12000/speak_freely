package models

type (
	File struct {
		ID  int64  `json:"id,omitempty"`
		Key string `json:"key,omitempty"`
	}

	FileCU struct {
		Key *string `json:"key,omitempty"`
	}
)
