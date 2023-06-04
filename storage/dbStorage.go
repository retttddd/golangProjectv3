package storage

import (
	"awesomeProject3/service"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Record struct {
	Key   string `db:"key_data"`
	Value string `db:"value_data"`
}

type dbStorage struct {
	db *sqlx.DB
}

func NewDbStorage(url string) *dbStorage {
	m, err := migrate.New(
		"file://./scripts/migration",
		url)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	dbConnection, err := sqlx.Connect("postgres", url)
	if err != nil {
		log.Fatalln(err)
	}
	return &dbStorage{
		db: dbConnection,
	}
}

func (st dbStorage) Read(key string) (*service.StorageModel, error) {
	r := Record{}
	err := st.db.Get(&r, "SELECT * FROM secret_data WHERE key_data=$1", key)
	if err != nil {
		return nil, errors.New("Database storage cant read data." + err.Error())
	}

	return &service.StorageModel{Value: &r.Value}, nil
}

func (st dbStorage) Write(key string, model *service.StorageModel) error {
	if model == nil || model.Value == nil {
		return errors.New("you did not pass any value")
	}

	_, err := st.db.Exec(`INSERT INTO secret_data (key_data,value_data) VALUES ($1,$2) ON CONFLICT (key_data) DO UPDATE SET value_data = $2`, key, *model.Value)
	if err != nil {
		return errors.New("Database storage cant write data." + err.Error())
	}

	return nil

}
