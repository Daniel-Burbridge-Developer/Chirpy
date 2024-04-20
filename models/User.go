package models

type User struct {
	Email    interface{} `json:"email"`
	Id       int         `json:"id"`
	Password string      `json:"password"`
}

type UserWithoutPassword struct {
	Email interface{} `json:"email"`
	Id    int         `json:"id"`
}

func (db *DB) CreateUser(email string, password string) (UserWithoutPassword, error) {
	users, err := db.GetUsers()
	userMap := make(map[int]User)
	if err != nil {
		return UserWithoutPassword{}, err
	}

	user := User{
		Id:       len(users) + 1,
		Email:    email,
		Password: password,
	}

	for _, usr := range users {
		userMap[usr.Id] = usr
	}

	userMap[user.Id] = user

	dbStructure := DBStructure{
		Users: userMap,
	}

	err = db.writeDB(dbStructure)
	if err != nil {
		return UserWithoutPassword{}, err
	}

	userWithoutPass := UserWithoutPassword{
		Email: user.Email,
		Id:    user.Id,
	}

	return userWithoutPass, nil
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

func (db *DB) GetUsersWithoutPassword() ([]UserWithoutPassword, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()

	if err != nil {
		return make([]UserWithoutPassword, 0), err
	}

	users := make([]UserWithoutPassword, 0)

	for _, user := range dbStructure.Users {
		users = append(users, UserWithoutPassword{
			Email: user.Email,
			Id:    user.Id,
		})
	}

	return users, nil
}
