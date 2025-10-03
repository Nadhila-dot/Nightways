package notebook

import (
    store "nadhi.dev/sarvar/fun/database"
    "nadhi.dev/sarvar/fun/db"
)

func DeleteNotebook(username string, id int) error {
    return store.DeleteNotebook(db.NotebooksDB, username, id)
}

func DeleteItemFromNotebook(username string, id int, itemName string) error {
    return store.DeleteItemFromNotebook(db.NotebooksDB, username, id, itemName)
}