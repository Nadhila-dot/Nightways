package store

import (
    "encoding/json"
    "os"
    "path/filepath"
)


func InitDB(name string) (*DB, error) {
    dbPath := filepath.Join("./zp-database/", name)
    err := os.MkdirAll(dbPath, 0755)
    if err != nil {
        return nil, err
    }
    return &DB{Path: dbPath}, nil
}

func (db *DB) GetStore(name string) (*Store, error) {
    storePath := filepath.Join(db.Path, name+".json")
    if _, err := os.Stat(storePath); os.IsNotExist(err) {
        file, err := os.Create(storePath)
        if err != nil {
            return nil, err
        }
        file.Write([]byte("{}"))
        file.Close()
    }
    return &Store{Name: name, Path: storePath}, nil
}

func (s *Store) GetData(out any) error {
    bytes, err := os.ReadFile(s.Path)
    if err != nil {
        return err
    }
    return json.Unmarshal(bytes, out)
}

func (s *Store) SetData(data any) error {
    bytes, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(s.Path, bytes, 0644)
}