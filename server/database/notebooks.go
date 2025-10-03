package store

import (
    "fmt"
    "time"
    "math/rand"
)

// CreateNotebook creates a new notebook for a user
func CreateNotebook(db *DB, username, name, description string, optional Optional) (*Notebook, error) {
    store, err := db.GetStore("notebooks")
    if err != nil {
        return nil, err
    }
    
    // Get all notebooks for structure
    var notebooks map[string]map[string]Notebook // username -> id -> notebook
    if err := store.GetData(&notebooks); err != nil {
        notebooks = make(map[string]map[string]Notebook)
    }
    
    // Create user's notebooks map if it doesn't exist
    if _, exists := notebooks[username]; !exists {
        notebooks[username] = make(map[string]Notebook)
    }
    
    // Generate a unique ID
    id := 1 + rand.Intn(99999)
    idStr := fmt.Sprintf("%d", id)
    
    // Ensure ID is unique
    for {
        if _, exists := notebooks[username][idStr]; !exists {
            break
        }
        id = 1 + rand.Intn(99999)
        idStr = fmt.Sprintf("%d", id)
    }
    
    now := time.Now()
    notebook := Notebook{
        ID:          id,
        Name:        name,
        Username:    username,
        Description: description,
        CreatedAt:   now,
        UpdatedAt:   now,
        Optional:    optional,
        Items:       make(map[string]string),
    }
    
    notebooks[username][idStr] = notebook
    
    if err := store.SetData(notebooks); err != nil {
        return nil, err
    }
    
    return &notebook, nil
}

// GetNotebook retrieves a notebook by ID
func GetNotebook(db *DB, username string, id int) (*Notebook, error) {
    store, err := db.GetStore("notebooks")
    if err != nil {
        return nil, err
    }
    
    var notebooks map[string]map[string]Notebook
    if err := store.GetData(&notebooks); err != nil {
        return nil, err
    }
    
    userNotebooks, exists := notebooks[username]
    if !exists {
        return nil, fmt.Errorf("no notebooks found for user %s", username)
    }
    
    idStr := fmt.Sprintf("%d", id)
    notebook, exists := userNotebooks[idStr]
    if !exists {
        return nil, fmt.Errorf("notebook %d not found", id)
    }
    
    return &notebook, nil
}

// GetAllNotebooks retrieves all notebooks for a user
func GetAllNotebooks(db *DB, username string) ([]Notebook, error) {
    store, err := db.GetStore("notebooks")
    if err != nil {
        return nil, err
    }
    
    var notebooks map[string]map[string]Notebook
    if err := store.GetData(&notebooks); err != nil {
        return []Notebook{}, nil // Return empty list if no notebooks exist
    }
    
    userNotebooks, exists := notebooks[username]
    if !exists {
        return []Notebook{}, nil // Return empty list if user has no notebooks
    }
    
    result := make([]Notebook, 0, len(userNotebooks))
    for _, notebook := range userNotebooks {
        result = append(result, notebook)
    }
    
    return result, nil
}

// AddItemToNotebook adds a sheet to a notebook
func AddItemToNotebook(db *DB, username string, id int, sheetName, url string) error {
    store, err := db.GetStore("notebooks")
    if err != nil {
        return err
    }
    
    var notebooks map[string]map[string]Notebook
    if err := store.GetData(&notebooks); err != nil {
        return err
    }
    
    userNotebooks, exists := notebooks[username]
    if !exists {
        return fmt.Errorf("no notebooks found for user %s", username)
    }
    
    idStr := fmt.Sprintf("%d", id)
    notebook, exists := userNotebooks[idStr]
    if !exists {
        return fmt.Errorf("notebook %d not found", id)
    }
    
    notebook.Items[sheetName] = url
    notebook.UpdatedAt = time.Now()
    userNotebooks[idStr] = notebook
    notebooks[username] = userNotebooks
    
    return store.SetData(notebooks)
}

// DeleteNotebook removes a notebook
func DeleteNotebook(db *DB, username string, id int) error {
    store, err := db.GetStore("notebooks")
    if err != nil {
        return err
    }
    
    var notebooks map[string]map[string]Notebook
    if err := store.GetData(&notebooks); err != nil {
        return err
    }
    
    userNotebooks, exists := notebooks[username]
    if !exists {
        return nil // Already gone, nothing to do
    }
    
    idStr := fmt.Sprintf("%d", id)
    if _, exists := userNotebooks[idStr]; !exists {
        return nil // Already gone, nothing to do
    }
    
    delete(userNotebooks, idStr)
    notebooks[username] = userNotebooks
    
    return store.SetData(notebooks)
}

// DeleteItemFromNotebook removes a sheet from a notebook
func DeleteItemFromNotebook(db *DB, username string, id int, itemName string) error {
    store, err := db.GetStore("notebooks")
    if err != nil {
        return err
    }
    
    var notebooks map[string]map[string]Notebook
    if err := store.GetData(&notebooks); err != nil {
        return err
    }
    
    userNotebooks, exists := notebooks[username]
    if !exists {
        return fmt.Errorf("no notebooks found for user %s", username)
    }
    
    idStr := fmt.Sprintf("%d", id)
    notebook, exists := userNotebooks[idStr]
    if !exists {
        return fmt.Errorf("notebook %d not found", id)
    }
    
    if _, exists := notebook.Items[itemName]; !exists {
        return fmt.Errorf("item %s not found in notebook", itemName)
    }
    
    delete(notebook.Items, itemName)
    notebook.UpdatedAt = time.Now()
    userNotebooks[idStr] = notebook
    notebooks[username] = userNotebooks
    
    return store.SetData(notebooks)
}

// GetItemsInNotebook gets all sheets in a notebook
func GetItemsInNotebook(db *DB, username string, id int) (map[string]string, error) {
    notebook, err := GetNotebook(db, username, id)
    if err != nil {
        return nil, err
    }
    
    return notebook.Items, nil
}

func UpdateNotebook(db *DB, username string, notebook Notebook) error {
    store, err := db.GetStore("notebooks")
    if err != nil {
        return err
    }
    
    var notebooks map[string]map[string]Notebook
    if err := store.GetData(&notebooks); err != nil {
        return err
    }
    
    userNotebooks, exists := notebooks[username]
    if !exists {
        return fmt.Errorf("no notebooks found for user %s", username)
    }
    
    idStr := fmt.Sprintf("%d", notebook.ID)
    if _, exists := userNotebooks[idStr]; !exists {
        return fmt.Errorf("notebook %d not found", notebook.ID)
    }
    
    // Keep same timestamps but update the rest
    originalNotebook := userNotebooks[idStr]
    notebook.CreatedAt = originalNotebook.CreatedAt
    notebook.UpdatedAt = time.Now()
    
    userNotebooks[idStr] = notebook
    notebooks[username] = userNotebooks
    
    return store.SetData(notebooks)
}