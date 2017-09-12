package models

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"
)

type UserSession struct {
	SessionKey   string    `json:"sessionKey"`
	UserID       int       `json:"userId"`
	LoginTime    time.Time `json:"loginTime"`
	LastSeenTime time.Time `json:"lastSeenTime"`
}

// NewUserSession given a userID returns a valid userSession stored in db
func NewUserSession(userID int) (string, error) {
	var k string
	key, err := makeKey()
	if err != nil {
		return "", err
	}

	t := time.Now()
	err = db.QueryRow(`INSERT into usersessions(sessionkey, userid, logintime, lastseentime)
						VALUES($1, $2, $3, $4)
						RETURNING sessionkey`, key, userID, t, t).Scan(&k)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return k, nil
}

func makeKey() (string, error) {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)

	return s, nil
}
