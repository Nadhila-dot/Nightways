package sheet

import (
    "fmt"
    _ "os"
    _ "sync"
	"log"
)


func (sq *SheetQueue) DeleteJob(id string) error {
    log.Printf("DeleteJob: locking queue")
    sq.mu.Lock()
    defer sq.mu.Unlock()
    log.Printf("DeleteJob: locked")

    // Stop job listener if exists
    if listener, ok := sq.jobListeners[id]; ok && listener != nil {
        delete(sq.jobListeners, id)
    }
    log.Printf("DeleteJob: loading jobs")
    jobs, err := sq.loadJobs()
    if err != nil {
        log.Printf("DeleteJob: failed to load jobs: %v", err)
        return fmt.Errorf("failed to load jobs: %w", err)
    }
    log.Printf("DeleteJob: loaded jobs")

    if _, exists := jobs[id]; !exists {
        log.Printf("DeleteJob: job %s not found", id)
        return fmt.Errorf("job %s not found", id)
    }
    delete(jobs, id)
    log.Printf("DeleteJob: saving jobs")
    if err := sq.saveJobs(jobs); err != nil {
        log.Printf("DeleteJob: failed to save jobs: %v", err)
        return fmt.Errorf("failed to save jobs: %w", err)
    }
    log.Printf("DeleteJob: done")
    return nil
}