package storages

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"sync"
)

type PostgresStorage struct {
	ConnString string
	pool       *pgxpool.Pool

	WaitGroup   sync.WaitGroup
	CloseSignal chan bool
	OsSignal    chan os.Signal
}

func (s *PostgresStorage) Connect() error {
	var err error
	s.pool, err = pgxpool.Connect(context.Background(), s.ConnString)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) Flush() error {

	return nil
}

func (s *PostgresStorage) Save(message []byte) error {
	return nil
}

func (s *PostgresStorage) Close() error {
	s.CloseSignal <- true

	s.WaitGroup.Wait()
	return nil
}

func (s *PostgresStorage) Start(messages chan LogMessage) {

	s.WaitGroup.Add(1)
	defer s.WaitGroup.Done()

	var msg LogMessage

	for {
		select {
		case msg = <-messages:
		//
		case <-s.CloseSignal:
			log.Println("CloseSignal received, closing...")
			return
		case sig := <-s.OsSignal:
			log.Printf("OsSignal received, closing...\n, %s\n", sig.String())
			return
		}
	}
}
