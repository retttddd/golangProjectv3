package storage

import (
	"errors"
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
	dbConnection, err := sqlx.Connect("postgres", url)
	if err != nil {
		log.Fatalln(err)
	}
	return &dbStorage{
		db: dbConnection,
	}
}

func (st dbStorage) Read(key string) (string, error) {
	r := Record{}
	err := st.db.Get(&r, "SELECT * FROM secret_data WHERE key_data=$1", key)
	if err != nil {
		return "", errors.New("Database storage cant read data." + err.Error())
	}

	return r.Value, nil
}

func (st dbStorage) Write(key string, value string) error {

	_, err := st.db.Exec(`INSERT INTO secret_data (key_data,value_data) VALUES ($1,$2) ON CONFLICT (key_data) DO UPDATE SET value_data = $2`, key, value)
	if err != nil {
		return errors.New("Database storage cant write data." + err.Error())
	}

	return nil

}
