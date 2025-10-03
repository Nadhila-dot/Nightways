package store

import (
    "time"
)

type QueuedJob struct {
    ID                string                 `json:"id"`
    UserID            string                 `json:"userId"`
    Prompt            string                 `json:"prompt"`
    Status            string                 `json:"status"`
    Result            interface{}            `json:"result,omitempty"`
    CreatedAt         time.Time              `json:"createdAt"`
    UpdatedAt         time.Time              `json:"updatedAt"`
    ConnectionCreated time.Time              `json:"connectionCreated"`
    Retries           int                    `json:"retries"`
    MaxRetry          int                    `json:"maxRetry"`
    Data              map[string]interface{} `json:"data,omitempty"`
}

// AddQueuedJob adds a job to the queue
func AddQueuedJob(db *DB, job QueuedJob) error {
    store, err := db.GetStore("queue")
    if err != nil {
        return err
    }
    
    var jobs map[string]QueuedJob
    if err := store.GetData(&jobs); err != nil {
        jobs = make(map[string]QueuedJob)
    }
    
    jobs[job.ID] = job
    return store.SetData(jobs)
}

// GetQueuedJob gets a job by ID
func GetQueuedJob(db *DB, id string) (*QueuedJob, error) {
    store, err := db.GetStore("queue")
    if err != nil {
        return nil, err
    }
    
    var jobs map[string]QueuedJob
    if err := store.GetData(&jobs); err != nil {
        return nil, err
    }
    
    job, ok := jobs[id]
    if !ok {
        return nil, nil
    }
    
    return &job, nil
}

// GetAllQueuedJobs gets all jobs
func GetAllQueuedJobs(db *DB, status string) (map[string]QueuedJob, error) {
    store, err := db.GetStore("queue")
    if err != nil {
        return nil, err
    }
    
    var jobs map[string]QueuedJob
    if err := store.GetData(&jobs); err != nil {
        jobs = make(map[string]QueuedJob)
    }
    
    if status == "" {
        return jobs, nil
    }
    
    // Filter by status if provided
    filteredJobs := make(map[string]QueuedJob)
    for id, job := range jobs {
        if job.Status == status {
            filteredJobs[id] = job
        }
    }
    
    return filteredJobs, nil
}

// UpdateQueuedJobStatus updates a job's status and result
func UpdateQueuedJobStatus(db *DB, id, status string, result interface{}) error {
    store, err := db.GetStore("queue")
    if err != nil {
        return err
    }
    
    var jobs map[string]QueuedJob
    if err := store.GetData(&jobs); err != nil {
        return err
    }
    
    job, ok := jobs[id]
    if !ok {
        return nil // Job not found, silently ignore
    }
    
    job.Status = status
    job.Result = result
    job.UpdatedAt = time.Now()
    jobs[id] = job
    
    return store.SetData(jobs)
}

// GetQueuedJobsByUser gets all jobs for a specific user
func GetQueuedJobsByUser(db *DB, userID string) ([]QueuedJob, error) {
    store, err := db.GetStore("queue")
    if err != nil {
        return nil, err
    }
    
    var jobs map[string]QueuedJob
    if err := store.GetData(&jobs); err != nil {
        return nil, err
    }
    
    var userJobs []QueuedJob
    for _, job := range jobs {
        if job.UserID == userID {
            userJobs = append(userJobs, job)
        }
    }
    
    return userJobs, nil
}

// RemoveQueuedJob removes a job from the queue
func RemoveQueuedJob(db *DB, id string) error {
    store, err := db.GetStore("queue")
    if err != nil {
        return err
    }
    
    var jobs map[string]QueuedJob
    if err := store.GetData(&jobs); err != nil {
        return err
    }
    
    delete(jobs, id)
    return store.SetData(jobs)
}

// CleanupOldJobs removes jobs older than the specified duration
func CleanupOldJobs(db *DB, maxAge time.Duration) error {
    store, err := db.GetStore("queue")
    if err != nil {
        return err
    }
    
    var jobs map[string]QueuedJob
    if err := store.GetData(&jobs); err != nil {
        return err
    }
    
    now := time.Now()
    for id, job := range jobs {
        if now.Sub(job.UpdatedAt) > maxAge {
            delete(jobs, id)
        }
    }
    
    return store.SetData(jobs)
}