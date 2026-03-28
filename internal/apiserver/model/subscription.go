package model

// Subscription представляет подписку
type Subscription struct {
	ID          int64  `json:"id" example:"1" description:"ID подписки"`
	ServiceName string `json:"service_name" example:"Netflix" description:"Название сервиса"`
	Price       int64  `json:"price" example:"999" description:"Стоимость подписки"`
	UserID      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba" description:"ID пользователя"`
	StartDate   string `json:"start_date" example:"03-2025" description:"Дата начала подписки"`
}
