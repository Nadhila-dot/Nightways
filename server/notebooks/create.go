package notebook

import (
    store "nadhi.dev/sarvar/fun/database"
    "nadhi.dev/sarvar/fun/db"
)

func CreateNotebook(username, name, description string, optional Optional) (*Notebook, error) {
    storeNotebook, err := store.CreateNotebook(db.NotebooksDB, username, name, description, store.Optional{
        Tags:        optional.Tags,
        Color:       optional.Color,
        Description: optional.Description,
    })
    if err != nil {
        return nil, err
    }
    
    return &Notebook{
        ID:          storeNotebook.ID,
        Name:        storeNotebook.Name,
        Description: storeNotebook.Description,
        Optional:    optional,
        Items:       storeNotebook.Items,
    }, nil
}

func CreateItemToNotebook(username string, id int, sheetName, url string) error {
    return store.AddItemToNotebook(db.NotebooksDB, username, id, sheetName, url)
}