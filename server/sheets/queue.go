package sheet

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	_ "sync"
	"time"

	"nadhi.dev/sarvar/fun/ai"
	_ "nadhi.dev/sarvar/fun/ai"
	store "nadhi.dev/sarvar/fun/database"
	"nadhi.dev/sarvar/fun/db"
	logg "nadhi.dev/sarvar/fun/logs"
	ws "nadhi.dev/sarvar/fun/websocket"
)

// QueuedJob represents a sheet generation job
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

// Convert between store.QueuedJob and sheet.QueuedJob
func toStoreJob(job QueuedJob) store.QueuedJob {
	return store.QueuedJob{
		ID:                job.ID,
		UserID:            job.UserID,
		Prompt:            job.Prompt,
		Status:            job.Status,
		Result:            job.Result,
		CreatedAt:         job.CreatedAt,
		UpdatedAt:         job.UpdatedAt,
		ConnectionCreated: job.ConnectionCreated,
		Retries:           job.Retries,
		MaxRetry:          job.MaxRetry,
		Data:              job.Data,
	}
}

func fromStoreJob(job store.QueuedJob) QueuedJob {
	return QueuedJob{
		ID:                job.ID,
		UserID:            job.UserID,
		Prompt:            job.Prompt,
		Status:            job.Status,
		Result:            job.Result,
		CreatedAt:         job.CreatedAt,
		UpdatedAt:         job.UpdatedAt,
		ConnectionCreated: job.ConnectionCreated,
		Retries:           job.Retries,
		MaxRetry:          job.MaxRetry,
		Data:              job.Data,
	}
}

var GlobalSheetGenerator *SheetGenerator

func NewSheetQueue(logger *log.Logger, queueDir string) (*SheetQueue, error) {
	if logger == nil {
		logger = log.Default()
	}

	// Ensure queue directory exists
	if err := os.MkdirAll(queueDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create queue directory: %w", err)
	}

	queueFile := filepath.Join(queueDir, "queue.json")

	// Check if queue file exists and is valid
	if _, err := os.Stat(queueFile); err == nil {
		// File exists, validate it
		file, err := os.ReadFile(queueFile)
		if err != nil {
			logger.Printf("Warning: Could not read queue file: %v", err)
		} else {
			var jobs map[string]QueuedJob
			if err := json.Unmarshal(file, &jobs); err != nil {
				logger.Printf("Warning: Corrupted queue file detected. Creating a new one.")
				// Reset the queue file
				if err := os.WriteFile(queueFile, []byte("{}"), 0644); err != nil {
					return nil, fmt.Errorf("failed to reset queue file: %w", err)
				}
			}
		}
	} else if os.IsNotExist(err) {
		// Create a new queue file
		if err := os.WriteFile(queueFile, []byte("{}"), 0644); err != nil {
			return nil, fmt.Errorf("failed to create queue file: %w", err)
		}
	} else {
		return nil, fmt.Errorf("failed to check queue file: %w", err)
	}

	return &SheetQueue{
		queueFile:     queueFile,
		statusUpdates: make(chan StatusUpdate, 100),
		logger:        logger,
		jobListeners:  make(map[string]func(StatusUpdate)),
	}, nil
}

func (sq *SheetQueue) loadJobs() (map[string]QueuedJob, error) {
	storeJobs, err := store.GetAllQueuedJobs(db.QueueDB, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load jobs: %w", err)
	}

	jobs := make(map[string]QueuedJob)
	for id, storeJob := range storeJobs {
		jobs[id] = fromStoreJob(storeJob)
	}

	return jobs, nil
}

func (sq *SheetQueue) saveJobs(jobs map[string]QueuedJob) error {
	// We can't directly save all jobs at once with our API
	// So we'll get current jobs and update them individually
	currentJobs, err := store.GetAllQueuedJobs(db.QueueDB, "")
	if err != nil {
		return fmt.Errorf("failed to get current jobs: %w", err)
	}

	// Delete any jobs that no longer exist
	for id := range currentJobs {
		if _, exists := jobs[id]; !exists {
			if err := store.RemoveQueuedJob(db.QueueDB, id); err != nil {
				return fmt.Errorf("failed to remove job %s: %w", id, err)
			}
		}
	}

	// Add or update jobs
	for id, job := range jobs {
		storeJob := toStoreJob(job)
		if err := store.AddQueuedJob(db.QueueDB, storeJob); err != nil {
			return fmt.Errorf("failed to save job %s: %w", id, err)
		}
	}

	return nil
}

// Start initializes worker goroutines for the queue
func (sq *SheetQueue) Start(ctx context.Context, workerCount int) {
	// sq.logger.Printf("Starting sheet queue with %d workers", workerCount)
	logg.Warning(fmt.Sprintf("Starting sheet queue with %d workers", workerCount))

	// Status handler goroutine
	sq.wg.Add(1)
	go sq.statusHandler(ctx)

	// Worker goroutines
	for i := 0; i < workerCount; i++ {
		sq.wg.Add(1)
		go sq.worker(ctx, i)
	}
}

// Stop gracefully shuts down the queue
func (sq *SheetQueue) Stop() {
	sq.logger.Println("Stopping sheet queue")
	close(sq.statusUpdates)
	sq.wg.Wait()
}

// EnqueueJob adds a new sheet generation job to the queue
func (sq *SheetQueue) EnqueueJob(id, userID, prompt string, maxRetry int) error {
	now := time.Now()
	job := QueuedJob{
		ID:                id,
		UserID:            userID,
		Prompt:            prompt,
		Status:            "queued",
		CreatedAt:         now,
		UpdatedAt:         now,
		ConnectionCreated: now,
		MaxRetry:          maxRetry,
		Data: map[string]interface{}{
			"processed":   false,
			"enqueued_at": now.Unix(),
		},
	}

	storeJob := toStoreJob(job)
	if err := store.AddQueuedJob(db.QueueDB, storeJob); err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	sq.logger.Printf("Enqueueing job %s for user %s", id, userID)
	return nil
}

// worker processes jobs from the queue
func (sq *SheetQueue) worker(ctx context.Context, id int) {
	defer sq.wg.Done()
	logg.Success(fmt.Sprintf("Worker %d started", id))

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			sq.logger.Printf("Worker %d shutting down due to context cancellation", id)
			return

		case <-ticker.C:
			// Try to claim a job atomically
			jobToProcess := sq.claimNextJob(id)
			if jobToProcess != nil {
				sq.processJob(ctx, jobToProcess)
			}
		}
	}
}

// claimNextJob atomically claims the next available queued job
func (sq *SheetQueue) claimNextJob(workerID int) *QueuedJob {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	// Load jobs from database
	jobs, err := sq.loadJobs()
	if err != nil {
		sq.logger.Printf("Worker %d: Error loading jobs: %v", workerID, err)
		return nil
	}

	// Find a queued job to process
	for jobID, job := range jobs {
		// Only process jobs that are truly queued
		if job.Status != "queued" {
			continue
		}

		// Check if already processed
		if job.Data != nil {
			if processed, ok := job.Data["processed"].(bool); ok && processed {
				sq.logger.Printf("Worker %d: Skipping already processed job %s", workerID, jobID)
				continue
			}
		}

		// Claim the job by updating status
		job.Status = "processing"
		job.UpdatedAt = time.Now()
		if job.Data == nil {
			job.Data = make(map[string]interface{})
		}
		job.Data["processing_started"] = time.Now().Unix()
		job.Data["worker_id"] = fmt.Sprintf("worker-%d", workerID)

		// Save immediately to claim it
		storeJob := toStoreJob(job)
		if err := store.AddQueuedJob(db.QueueDB, storeJob); err != nil {
			sq.logger.Printf("Worker %d: Failed to claim job %s: %v", workerID, jobID, err)
			continue
		}

		sq.logger.Printf("Worker %d: Successfully claimed job %s", workerID, jobID)
		jobCopy := job
		jobCopy.ID = jobID
		return &jobCopy
	}

	return nil
}

// processJob handles the actual sheet generation logic
func (sq *SheetQueue) processJob(ctx context.Context, job *QueuedJob) {
	sq.logger.Printf("Processing job %s", job.ID)

	// CRITICAL: Check if already processed before starting
	if job.Data != nil {
		if processed, ok := job.Data["processed"].(bool); ok && processed {
			sq.logger.Printf("Job %s already processed, skipping", job.ID)
			return
		}
	}

	// Job is already marked as processing by worker, so we can proceed directly

	// Send an initial "processing" status update
	update := ws.Start("Starting job processing", map[string]interface{}{})

	sq.statusUpdates <- StatusUpdate{
		ID:     job.ID,
		Status: "processing",
		Data:   update["data"].(map[string]interface{}),
	}

	// Create object in database/storage
	err := sq.createSheetObject(job)
	if err != nil {
		sq.handleJobError(job, err)
		return
	}

	update_create := ws.Stage("AI Generation", "Sending data to AI", map[string]interface{}{})

	sq.statusUpdates <- StatusUpdate{
		ID:     job.ID,
		Status: "processing-create",
		Data:   update_create["data"].(map[string]interface{}),
	}

	// Call AI generation
	result, err := sq.generateWithAI(job)
	if err != nil {
		sq.handleJobError(job, err)
		return
	}

	// Mark as processed
	jobs, _ := sq.loadJobs()
	if finalJob, exists := jobs[job.ID]; exists {
		if finalJob.Data == nil {
			finalJob.Data = make(map[string]interface{})
		}
		finalJob.Data["processed"] = true
		finalJob.Data["completed_at"] = time.Now().Unix()
		jobs[job.ID] = finalJob
		sq.saveJobs(jobs)
	}

	fmt.Println(result)
}

func (sq *SheetQueue) createSheetObject(job *QueuedJob) error {
	// It already stores in the queue DB, so we just log here
	// maybe later we can add more complex logic
	//sq.logger.Printf("Creating sheet object for job %s", job.ID)
	logg.Info(fmt.Sprintf("Creating sheet object for job %s", job.ID))

	return nil
}

// handleJobError processes failures with retry logic
func (sq *SheetQueue) handleJobError(job *QueuedJob, err error) {
	sq.logger.Printf("Error processing job %s: %v", job.ID, err)

	jobs, loadErr := sq.loadJobs()
	if loadErr != nil {
		sq.logger.Printf("Failed to load jobs for error handling: %v", loadErr)
		return
	}

	currentJob, exists := jobs[job.ID]
	if !exists {
		sq.logger.Printf("Job %s no longer exists", job.ID)
		return
	}

	currentJob.Retries++
	var status string
	if currentJob.Retries <= currentJob.MaxRetry {
		sq.logger.Printf("Retrying job %s (%d/%d)", job.ID, currentJob.Retries, currentJob.MaxRetry)
		status = "retrying"
	} else {
		sq.logger.Printf("Job %s failed after %d retries", job.ID, currentJob.Retries)
		status = "failed"
	}

	currentJob.Status = status
	currentJob.UpdatedAt = time.Now()
	jobs[job.ID] = currentJob

	if err := sq.saveJobs(jobs); err != nil {
		sq.logger.Printf("Failed to save jobs after error handling: %v", err)
	}

	// Send only one status update using the Error helper
	errorUpdate := ws.Retry(
		"Job failed, retrying...",
		map[string]interface{}{
			"retries":   currentJob.Retries,
			"maxRetry":  currentJob.MaxRetry,
			"willRetry": currentJob.Retries <= currentJob.MaxRetry,
		},
	)
	sq.statusUpdates <- StatusUpdate{
		ID:     job.ID,
		Status: status,
		Result: err.Error(),
		Data:   errorUpdate["data"].(map[string]interface{}),
	}
}

// statusHandler processes status updates
func (sq *SheetQueue) statusHandler(ctx context.Context) {
	defer sq.wg.Done()
	// sq.logger.Println("Status handler started")
	logg.Info("Status handler started")

	for {
		select {
		case <-ctx.Done():
			sq.logger.Println("Status handler shutting down")
			return

		case update, ok := <-sq.statusUpdates:
			if !ok {
				sq.logger.Println("Status handler shutting down due to closed channel")
				return
			}

			jobs, err := sq.loadJobs()
			if err != nil {
				sq.logger.Printf("Failed to load jobs for status update: %v", err)
				continue
			}

			if job, exists := jobs[update.ID]; exists {
				job.Status = update.Status
				job.Result = update.Result
				job.UpdatedAt = time.Now()
				jobs[update.ID] = job

				if err := sq.saveJobs(jobs); err != nil {
					sq.logger.Printf("Failed to save jobs after status update: %v", err)
				}

				// Notify job-specific listener
				sq.mu.Lock()
				if listener, ok := sq.jobListeners[update.ID]; ok && listener != nil {
					listener(update)
				}
				sq.mu.Unlock()

				sq.logger.Printf("Job %s status updated to %s (ready for websocket notification)", update.ID, update.Status)
			}
		}
	}
}

// updateJobStatus sends a status update to the status handler
// This is deprecated in favor of sending StatusUpdate directly to the channel
func (sq *SheetQueue) updateJobStatus(id, status string, result interface{}) {
	sq.statusUpdates <- StatusUpdate{
		ID:     id,
		Status: status,
		Result: result,
	}
}

// RegisterJobListener allows registering a callback for job status updates
func (sq *SheetQueue) RegisterJobListener(jobID string, cb func(StatusUpdate)) {
	sq.mu.Lock()
	defer sq.mu.Unlock()
	sq.jobListeners[jobID] = cb
}

// GetJobsByUser returns all jobs for a specific user
func (sq *SheetQueue) GetJobsByUser(userID string) ([]QueuedJob, error) {
	jobs, err := sq.loadJobs()
	if err != nil {
		return nil, fmt.Errorf("failed to load jobs: %w", err)
	}

	var userJobs []QueuedJob
	for _, job := range jobs {
		if job.UserID == userID {
			userJobs = append(userJobs, job)
		}
	}

	return userJobs, nil
}

// CleanupOldJobs removes jobs older than the specified duration
func (sq *SheetQueue) CleanupOldJobs(maxAge time.Duration) error {
	jobs, err := sq.loadJobs()
	if err != nil {
		return fmt.Errorf("failed to load jobs: %w", err)
	}

	now := time.Now()
	for id, job := range jobs {
		if now.Sub(job.UpdatedAt) > maxAge {
			delete(jobs, id)
		}
	}

	return sq.saveJobs(jobs)
}

func (sq *SheetQueue) GetJobStatus(id string) (QueuedJob, bool) {
	storeJob, err := store.GetQueuedJob(db.QueueDB, id)
	if err != nil || storeJob == nil {
		return QueuedJob{}, false
	}
	return fromStoreJob(*storeJob), true
}

// extractGenerationRequest extracts the GenerationRequest from a job's prompt
func (sq *SheetQueue) extractGenerationRequest(job *QueuedJob) (*ai.GenerationRequest, error) {
	// The prompt contains JSON with nested request
	var metadata struct {
		Request   json.RawMessage `json:"request"`
		Processed bool            `json:"processed"`
		CreatedAt int64           `json:"created_at"`
	}

	if err := json.Unmarshal([]byte(job.Prompt), &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job metadata: %w", err)
	}

	var request ai.GenerationRequest
	if err := json.Unmarshal(metadata.Request, &request); err != nil {
		return nil, fmt.Errorf("failed to unmarshal generation request: %w", err)
	}

	return &request, nil
}
