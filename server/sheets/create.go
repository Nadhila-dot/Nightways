package sheet

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"nadhi.dev/sarvar/fun/ai"
)

// SheetGenerator manages the sheet creation process
type SheetGenerator struct {
	Queue  *SheetQueue // Change from queue to Queue (uppercase for export)
	logger *log.Logger
}

// NewSheetGenerator creates a new sheet generator with queue system
func NewSheetGenerator(logger *log.Logger, queueDir string, workerCount int) (*SheetGenerator, error) {
	if logger == nil {
		logger = log.Default()
	}

	queue, err := NewSheetQueue(logger, queueDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue: %w", err)
	}

	generator := &SheetGenerator{
		Queue:  queue, // Change from queue to Queue
		logger: logger,
	}

	// Start the queue with background context
	queue.Start(context.Background(), workerCount)

	return generator, nil
}

func (sg *SheetGenerator) CreateSheet(userID string, request *ai.GenerationRequest) (string, error) {
	jobID := fmt.Sprintf("sheet-%d", time.Now().UnixNano())

	// Store the request directly as JSON in the Prompt field
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Use Data field for processing flags, Prompt for the actual request
	err = sg.Queue.EnqueueJob(jobID, userID, string(requestJSON), 3)
	if err != nil {
		return "", fmt.Errorf("failed to enqueue job: %w", err)
	}

	return jobID, nil
}

// GetUserJobs returns all jobs for a specific user
func (sg *SheetGenerator) GetUserJobs(userID string) ([]QueuedJob, error) {
	return sg.Queue.GetJobsByUser(userID)
}

// Shutdown stops the sheet generator and its queue
func (sg *SheetGenerator) Shutdown() {
	sg.Queue.Stop()
}

// CleanupOldJobs removes jobs older than the specified duration
func (sg *SheetGenerator) CleanupOldJobs(maxAge time.Duration) error {
	return sg.Queue.CleanupOldJobs(maxAge)
}
