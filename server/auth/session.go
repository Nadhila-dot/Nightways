package auth

import (
	"errors"
	"math/rand"
	"time"

	store "nadhi.dev/sarvar/fun/database"
	"nadhi.dev/sarvar/fun/db"
)

type Session struct {
    ID       string `json:"id"`
    Username string `json:"username"`
}

func generateSessionID() string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, 16)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return "nadhi.dev_" + string(b)
}

func GetUserBySession(sessionID string) (*store.User, error) {
    s, err := store.GetSession(db.SessionsDB, sessionID)
    if err != nil {
        return nil, err
    }
    if s == nil {
        return nil, errors.New("invalid session")
    }
    username, ok := s.Data["username"].(string)
    if !ok {
        return nil, errors.New("username not found")
    }
    users, err := loadUsers()
    if err != nil {
        return nil, err
    }
    for _, u := range users {
        if u.Username == username {
            return &u, nil
        }
    }
    return nil, errors.New("user not found")
}

func CreateSession(username string) (string, error) {
    id := generateSessionID()
    session := store.Session{ID: id, Data: map[string]any{"username": username}}
    err := store.AddSession(db.SessionsDB, session)
    if err != nil {
        return "", err
    }
    return id, nil
}

func IsSessionValid(sessionID string) (bool, error) {
    s, err := store.GetSession(db.SessionsDB, sessionID)
    if err != nil {
        return false, err
    }
    return s != nil, nil
}