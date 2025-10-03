package auth

import (
	"errors"
	"strings"

	store "nadhi.dev/sarvar/fun/database"
	"nadhi.dev/sarvar/fun/db"
)

func Register(username, email, password, rank string) error {
    users, err := store.GetAllUsers(db.UsersDB)
    if err != nil {
        return err
    }
    for _, u := range users {
        if strings.EqualFold(u.Username, username) || strings.EqualFold(u.Email, email) {
            return errors.New("username or email already exists")
        }
    }
    user := store.User{Username: username, Email: email, Password: password, Rank: rank}
    return store.AddUser(db.UsersDB, user)
}

func Login(identifier, password string) (string, error) {
    users, err := store.GetAllUsers(db.UsersDB)
    if err != nil {
        return "", err
    }
    for _, u := range users {
        if (strings.EqualFold(u.Username, identifier) || strings.EqualFold(u.Email, identifier)) && u.Password == password {
            return CreateSession(u.Username)
        }
    }
    return "", errors.New("invalid credentials")
}

func loadUsers() ([]store.User, error) {
    usersMap, err := store.GetAllUsers(db.UsersDB)
    if err != nil {
        return nil, err
    }
    var users []store.User
    for _, u := range usersMap {
        users = append(users, u)
    }
    return users, nil
}