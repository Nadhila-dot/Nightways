package store

type User struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Rank     string `json:"rank"`
}

func AddUser(db *DB, user User) error {
    store, err := db.GetStore("users")
    if err != nil {
        return err
    }
    var users map[string]User
    if err := store.GetData(&users); err != nil {
        users = make(map[string]User)
    }
    users[user.Username] = user
    return store.SetData(users)
}

func GetUser(db *DB, username string) (*User, error) {
    store, err := db.GetStore("users")
    if err != nil {
        return nil, err
    }
    var users map[string]User
    if err := store.GetData(&users); err != nil {
        return nil, err
    }
    user, ok := users[username]
    if !ok {
        return nil, nil
    }
    return &user, nil
}

func GetAllUsers(db *DB) (map[string]User, error) {
    store, err := db.GetStore("users")
    if err != nil {
        return nil, err
    }
    var users map[string]User
    if err := store.GetData(&users); err != nil {
        users = make(map[string]User)
    }
    return users, nil
}

func RemoveUser(db *DB, username string) error {
    store, err := db.GetStore("users")
    if err != nil {
        return err
    }
    var users map[string]User
    if err := store.GetData(&users); err != nil {
        return err
    }
    delete(users, username)
    return store.SetData(users)
}