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
func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	_, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := db.ensureDB()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	chirps, err := db.GetChirps()
	chirpMap := make(map[int]Chirp)
	if err != nil {
		return Chirp{}, err
	}
	chirp := Chirp{
		Id:   len(chirps) + 1, // Assuming chirp IDs start from 1
		Body: body,
	}

	// Update the map with chirp ID as key
	chirpMap[chirp.Id] = chirp

	dbStructure := DBStructure{
		Chirps: chirpMap,
	}

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

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
	// Create the file with empty content if it doesn't exist
	file, err := os.Create(db.path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Set permissions for the file
	err = file.Chmod(0666)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	var dbStructure DBStructure
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

// Current known Errors:
// not saving to file
// can't call get chirps on empty DB
