package vela

import (
    "sort"
    "strings"
    "time"

    store "nadhi.dev/sarvar/fun/database"
    "nadhi.dev/sarvar/fun/db"
)

// SheetQueueItem represents a single queue item
type SheetQueueItem struct {
    ID        string      `json:"id"`
    UserID    string      `json:"user_id"`
    Prompt    string      `json:"prompt"`
    Status    string      `json:"status"`
    CreatedAt string      `json:"created_at"`
    UpdatedAt string      `json:"updated_at"`
    Retries   int         `json:"retries"`
    MaxRetry  int         `json:"max_retry"`
    Result    interface{} `json:"result"`
}

// Convert store.QueuedJob to SheetQueueItem
func fromStoreJob(job store.QueuedJob) SheetQueueItem {
    return SheetQueueItem{
        ID:        job.ID,
        UserID:    job.UserID,
        Prompt:    job.Prompt,
        Status:    job.Status,
        CreatedAt: job.CreatedAt.Format(time.RFC3339), // Convert time.Time to string
        UpdatedAt: job.UpdatedAt.Format(time.RFC3339),
        Retries:   job.Retries,
        MaxRetry:  job.MaxRetry,
        Result:    job.Result,
    }
}

// GetQueueItems fetches queue items with pagination, sorting, and optional search
func GetQueueItems(queuePath string, latest bool, objNum int, search string) ([]SheetQueueItem, error) {
    // Get all jobs from the database (ignoring queuePath since we're using DB now)
    jobs, err := store.GetAllQueuedJobs(db.QueueDB, "")
    if err != nil {
        return nil, err
    }

    // Convert map to slice
    items := make([]SheetQueueItem, 0, len(jobs))
    for _, job := range jobs {
        items = append(items, fromStoreJob(job))
    }

    // Optional search filter (case-insensitive, checks ID, UserID, Prompt, Status)
    if search != "" {
        searchLower := strings.ToLower(search)
        filtered := make([]SheetQueueItem, 0, len(items))
        for _, item := range items {
            if strings.Contains(strings.ToLower(item.ID), searchLower) ||
                strings.Contains(strings.ToLower(item.UserID), searchLower) ||
                strings.Contains(strings.ToLower(item.Prompt), searchLower) ||
                strings.Contains(strings.ToLower(item.Status), searchLower) {
                filtered = append(filtered, item)
            }
        }
        items = filtered
    }

    // Sort by CreatedAt (latest first if requested)
    // Parse the string back to time for sorting
    sort.Slice(items, func(i, j int) bool {
        timeI, errI := time.Parse(time.RFC3339, items[i].CreatedAt)
        timeJ, errJ := time.Parse(time.RFC3339, items[j].CreatedAt)
        if errI != nil || errJ != nil {
            // If parsing fails, fall back to string comparison
            if latest {
                return items[i].CreatedAt > items[j].CreatedAt
            }
            return items[i].CreatedAt < items[j].CreatedAt
        }
        if latest {
            return timeI.After(timeJ)
        }
        return timeI.Before(timeJ)
    })

    // Pagination: return up to objNum items
    if objNum > 0 && objNum < len(items) {
        items = items[:objNum]
    }

    return items, nil
}