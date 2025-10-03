package notebook

import (
    store "nadhi.dev/sarvar/fun/database"
    "nadhi.dev/sarvar/fun/db"
)


func GetAllNotebooks(username string) ([]store.Notebook, error) {
    return store.GetAllNotebooks(db.NotebooksDB, username)
}


func GetNotebook(username string, id int) (*store.Notebook, error) {
    return store.GetNotebook(db.NotebooksDB, username, id)
}


func GetItemsInNotebook(username string, id int) (map[string]string, error) {
    return store.GetItemsInNotebook(db.NotebooksDB, username, id)
}