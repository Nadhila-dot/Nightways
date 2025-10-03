package sheet

import (
	_"time"
	"sync"
	"log"
)



type StatusUpdate struct {
	ID     string                 `json:"id"`
	Status string                 `json:"status"`
	Result interface{}            `json:"result,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"` // optional due to omitempty
}

type SheetQueue struct {
	queueFile     string
	statusUpdates chan StatusUpdate
	wg            sync.WaitGroup
	logger        *log.Logger
	mu            sync.Mutex
	jobListeners  map[string]func(StatusUpdate)
}