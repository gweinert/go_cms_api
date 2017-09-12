package models

import "log"
import "strings"

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	PasswordSalt string `json:"passwordSalt"`
	IsDisabled   bool   `json:"isDisabled"`
	Level        int    `json:"level"`
}

// LoginUser gets user by unique email. if password matches returns user and user sessionkey
func LoginUser(email string, passwordHash string) (*User, string, error) {
	user, err := FindUserByEmail(email)
	if err != nil {
		log.Fatal(err)
		return nil, "", err
	}

	if user.PasswordHash == passwordHash {
		sessionID, err := NewUserSession(user.ID)
		if err != nil {
			log.Fatal(err)
			return nil, "", err
		}
		return user, sessionID, nil
	}

	return nil, "", nil
}

func FindUserByEmail(email string) (*User, error) {
	user := new(User)
	err := db.QueryRow(`SELECT * from users WHERE email = $1`,
		email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.PasswordSalt, &user.IsDisabled, &user.Level)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return user, nil
}

func FindUserByID(id int) (*User, error) {
	user, err := findUserBy("id", id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return user, nil
}

func findUserBy(param string, paramValue int) (*User, error) {
	user := new(User)
	queryString := strings.Join([]string{"SELECT * from users WHERE ", param, " = $1"}, "")
	err := db.QueryRow(queryString,
		paramValue).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.PasswordSalt, &user.IsDisabled, &user.Level)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return user, nil
}

func GetUserFromSessionID(sessionID string) (*User, error) {
	var userID int

	err := db.QueryRow(`SELECT userid from usersessions 
						WHERE sessionkey = $1`, sessionID).Scan(&userID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	user, err := FindUserByID(userID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return user, nil
}
