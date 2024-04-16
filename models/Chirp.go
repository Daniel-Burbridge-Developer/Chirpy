package models

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

// Only used in writing, maybe reading. I don't need one of these because I'm writing and reading every update. remember this.
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int         `json:"id"`
	Body interface{} `json:"body"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func (db *DB) NewDB(path string) (*DB, error) {
	_, err := os.ReadFile(path)
	if err != nil {
		if err == os.ErrNotExist {
			db.ensureDB()
		}
		return nil, err
	}
	return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	chirps, err := db.GetChirps()
	if err != nil {
		return Chirp{}, err
	}
	chirp := Chirp{
		Id:   len(chirps) + 1,
		Body: body,
	}

	chirps = append(chirps, chirp)

	dbStructure := DBStructure{}

	for i, chirp := range chirps {
		dbStructure.Chirps[i] = chirp
	}

	db.writeDB(dbStructure)

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()

	if err != nil {
		return make([]Chirp, 0), err
	}

	chirps := make([]Chirp, 0)

	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	os.WriteFile(db.path, []byte(""), 0666)

	return nil
}

// loadDB reads the database file int0o memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	JSON, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, JSON, 0666)
	if err != nil {
		return err
	}

	return nil
}