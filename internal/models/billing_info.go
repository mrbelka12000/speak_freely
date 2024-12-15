package models

import "time"

type (
	BillingInfo struct {
		ID        int64
		UserID    int64 // telegram UID
		ChatID    int64
		DebitDate time.Time
	}

	BillingInfoCU struct {
		UserID    *int64
		ChatID    *int64
		DebitDate *time.Time
	}
)
