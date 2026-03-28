package store

import (
	"apiserver/internal/apiserver/model"
	"database/sql"
)

type SubscriptionRepository struct {
	store *Store
}

func (s *SubscriptionRepository) Get(id int64) (*model.Subscription, error) {
	var sub model.Subscription
	if err := s.store.db.QueryRow("SELECT id,service_name,price,user_id,TO_CHAR(start_date,'DD-MM-YYYY') FROM subscription WHERE id = $1", id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
	); err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *SubscriptionRepository) Create(subs *model.Subscription) error {
	if err := s.store.db.QueryRow(
		"INSERT INTO subscription(service_name,price,user_id,start_date) VALUES($1,$2,$3,$4) RETURNING id",
		subs.ServiceName,
		subs.Price,
		subs.UserID,
		subs.StartDate,
	).Scan(&subs.ID); err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionRepository) Delete(id int64) error {
	if _, err := s.store.db.Exec("DELETE FROM subscription WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionRepository) Update(subs *model.Subscription) error {
	if _, err := s.store.db.Exec(
		"UPDATE subscription SET service_name = $1 ,price = $2,user_id = $3,start_date = $4 WHERE id = $5",
		subs.ServiceName,
		subs.Price,
		subs.UserID,
		subs.StartDate,
		subs.ID,
	); err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionRepository) List() ([]model.Subscription, error) {
	rows, err := s.store.db.Query("SELECT id,service_name,price,user_id,TO_CHAR(start_date,'DD-MM-YYYY') FROM subscription")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	list := make([]model.Subscription, 0)
	for rows.Next() {
		var row model.Subscription
		err := rows.Scan(
			&row.ID,
			&row.ServiceName,
			&row.Price,
			&row.UserID,
			&row.StartDate,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, row)
	}
	return list, nil
}

func (s *SubscriptionRepository) SummaryByIdAndPeriod(uuid, service_name, start_date, end_date string) (int64, error) {
	var price sql.NullInt64

	if err := s.store.db.QueryRow(
		"SELECT SUM(price) FROM subscription WHERE user_id = $1 AND service_name = $2 AND start_date BETWEEN $3 AND $4",
		uuid, service_name, start_date, end_date,
	).Scan(&price); err != nil {
		return 0, err
	}
	return int64(price.Int64), nil
}
