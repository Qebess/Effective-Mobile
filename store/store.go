package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	db                     *sql.DB
	subscriptionRepository *SubscriptionRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) Open(dbpath string) error {
	db, err := sql.Open("postgres", dbpath)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) User() *SubscriptionRepository {
	if s.subscriptionRepository != nil {
		return s.subscriptionRepository
	}

	s.subscriptionRepository = &SubscriptionRepository{
		store: s,
	}
	return s.subscriptionRepository

}
