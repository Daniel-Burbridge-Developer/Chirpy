package models

import (
	"encoding/json"
	"fmt"
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
	Users  map[int]User  `json:"users"`
}
type Chirp struct {
	Body interface{} `json:"body"`
	Id   int         `json:"id"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist

// CREATING A USER DELETES ALL CHIRPS
// CREATING A CHIRP DELETES ALL USERS
// I DON'T KNOW WHY

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
	// fmt.Printf("inside createchirp %v\n", body)
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
	for _, ch := range chirps {
		chirpMap[ch.Id] = ch
	}

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

	// fmt.Printf("file path %v\n", db.path)
	data, err := os.ReadFile(db.path)
	if err != nil {
		fmt.Printf("unable to read file?")
		return DBStructure{}, err
	}

	// fmt.Println("Data read from file:", string(data))

	dbStructure := DBStructure{}
	err = json.Unmarshal(data, &dbStructure)

	if err != nil {
		// fmt.Printf("unable to Unmarshal file? %v\n", err)
		return DBStructure{}, nil
	}

	// fmt.Printf("dbs not doing what I think it does?\n")
	// fmt.Printf("%v", dbStructure.Chirps)
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
