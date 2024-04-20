package models

type User struct {
	Email interface{} `json:"email"`
	Id    int         `json:"id"`
}

func (db *DB) CreateUser(email string) (User, error) {
	// fmt.Printf("inside createchirp %v\n", body)
	users, err := db.GetUsers()
	userMap := make(map[int]User)
	if err != nil {
		return User{}, err
	}
	user := User{
		Id:    len(users) + 1, // Assuming chirp IDs start from 1
		Email: email,
	}

	// Update the map with chirp ID as key
	for _, usr := range users {
		userMap[usr.Id] = usr
	}

	userMap[user.Id] = user

	dbStructure := DBStructure{
		Users: userMap,
	}

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetUsers() ([]User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()

	if err != nil {
		return make([]User, 0), err
	}

	users := make([]User, 0)

	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}
