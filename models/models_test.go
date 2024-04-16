package models

import (
	"os"
	"sync"
	"testing"
)

func TestWriteDB(t *testing.T) {
	// Create a temporary file for testing
	filePath := "test_db.json"

	// Create a new DB instance
	db := &DB{
		path: filePath,
		mux:  &sync.RWMutex{},
	}

	// Define test data
	testChirps := []Chirp{
		{Id: 1, Body: "Hello, world!"},
		{Id: 2, Body: "Testing chirps"},
	}

	testDBStructure := DBStructure{
		Chirps: make(map[int]Chirp),
	}

	for i, chirp := range testChirps {
		testDBStructure.Chirps[i] = chirp
	}

	// Write test data to the file
	err := db.writeDB(testDBStructure)
	if err != nil {
		t.Errorf("writeDB failed: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("writeDB did not create the file")
	}

	// Clean up: remove the temporary file
	err = os.Remove(filePath)
	if err != nil {
		t.Errorf("cleanup failed: %v", err)
	}
}
