package store


type Session struct {
    ID   string
    Data map[string]any
}

func AddSession(db *DB, session Session) error {
    store, err := db.GetStore("sessions")
    if err != nil {
        return err
    }
    var sessions map[string]Session
    if err := store.GetData(&sessions); err != nil {
        sessions = make(map[string]Session)
    }
    sessions[session.ID] = session
    return store.SetData(sessions)
}

func GetSession(db *DB, id string) (*Session, error) {
    store, err := db.GetStore("sessions")
    if err != nil {
        return nil, err
    }
    var sessions map[string]Session
    if err := store.GetData(&sessions); err != nil {
        return nil, err
    }
    session, ok := sessions[id]
    if !ok {
        return nil, nil
    }
    return &session, nil
}

func RemoveSession(db *DB, id string) error {
    store, err := db.GetStore("sessions")
    if err != nil {
        return err
    }
    var sessions map[string]Session
    if err := store.GetData(&sessions); err != nil {
        return err
    }
    delete(sessions, id)
    return store.SetData(sessions)
}