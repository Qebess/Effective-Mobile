package model

type Subscription struct {
	ID          int64  `json:"id,omitempty"`
	ServiceName string `json:"service_name"`
	Price       int64  `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
}
