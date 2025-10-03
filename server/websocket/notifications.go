package websocket

import (
    "encoding/json"
    "log"
    "sync"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
    "nadhi.dev/sarvar/fun/auth"
)

// Connection represents a websocket connection with user info
type Connection struct {
    Conn      *websocket.Conn
    UserID    string
    Connected time.Time
}

// UserSocketManager handles websocket connections and messaging
type UserSocketManager struct {
    connections map[string][]*Connection // map[userID][]Connection
    mu          sync.RWMutex
    logger      *log.Logger
}

// NewUserSocketManager creates a new UserSocketManager
func NewUserSocketManager(logger *log.Logger) *UserSocketManager {
    if logger == nil {
        logger = log.Default()
    }
    
    return &UserSocketManager{
        connections: make(map[string][]*Connection),
        logger:      logger,
    }
}

// ConnectHandler handles new websocket connections
func (usm *UserSocketManager) ConnectHandler(c *websocket.Conn) {
    sessionID := c.Query("session")
    if sessionID == "" {
        usm.logger.Println("Websocket connection attempt without session ID")
        c.Close()
        return
    }

    // Validate session
    valid, err := auth.IsSessionValid(sessionID)
    if err != nil {
        usm.logger.Printf("Error validating session: %v", err)
        c.Close()
        return
    }

    if !valid {
        usm.logger.Printf("Invalid session ID: %s", sessionID)
        c.Close()
        return
    }

    // Get user info from session
    user, err := auth.GetUserBySession(sessionID)
    if err != nil {
        usm.logger.Printf("Error getting user from session: %v", err)
        c.Close()
        return
    }

    // Create connection
    conn := &Connection{
        Conn:      c,
        UserID:    user.Username,
        Connected: time.Now(),
    }

    // Add to connections map
    usm.mu.Lock()
    if _, exists := usm.connections[user.Username]; !exists {
        usm.connections[user.Username] = make([]*Connection, 0)
    }
    usm.connections[user.Username] = append(usm.connections[user.Username], conn)
    usm.mu.Unlock()

    usm.logger.Printf("User %s connected via websocket", user.Username)

    // Send welcome message
    welcomeMsg := map[string]interface{}{
        "type":    "welcome",
        "message": "Connected to Vela notification system",
        "time":    time.Now(),
    }
    
    err = conn.Conn.WriteJSON(welcomeMsg)
    if err != nil {
        usm.logger.Printf("Error sending welcome message: %v", err)
    }

    // Handle incoming messages
    for {
        msgType, msg, err := c.ReadMessage()
        if err != nil {
            usm.logger.Printf("Error reading message: %v", err)
            break
        }

        // Log incoming messages
        usm.logger.Printf("Received message from %s: %s", user.Username, string(msg))

        // Echo messages back (for testing)
        if msgType == websocket.TextMessage {
            var data map[string]interface{}
            if err := json.Unmarshal(msg, &data); err != nil {
                usm.logger.Printf("Error parsing message: %v", err)
                continue
            }

            // Add timestamp and echo
            data["received"] = time.Now()
            data["echo"] = true
            
            responseMsg, err := json.Marshal(data)
            if err != nil {
                usm.logger.Printf("Error creating response: %v", err)
                continue
            }
            
            if err := c.WriteMessage(websocket.TextMessage, responseMsg); err != nil {
                usm.logger.Printf("Error sending echo response: %v", err)
            }
        }
    }

    // Handle disconnection
    usm.mu.Lock()
    defer usm.mu.Unlock()

    // Remove connection from slice
    conns := usm.connections[user.Username]
    for i, connection := range conns {
        if connection.Conn == c {
            usm.connections[user.Username] = append(conns[:i], conns[i+1:]...)
            break
        }
    }

    // If no connections left for user, remove the user entry
    if len(usm.connections[user.Username]) == 0 {
        delete(usm.connections, user.Username)
    }

    usm.logger.Printf("User %s disconnected from websocket", user.Username)
}

// Send sends a message to a specific user
func (usm *UserSocketManager) Send(userID string, data interface{}) error {
    usm.mu.RLock()
    defer usm.mu.RUnlock()

    conns, exists := usm.connections[userID]
    if !exists || len(conns) == 0 {
        return fiber.NewError(fiber.StatusNotFound, "User not connected")
    }

    msg := map[string]interface{}{
        "type":      "notification",
        "timestamp": time.Now(),
        "data":      data,
    }

    // Log what we're sending
    msgBytes, _ := json.Marshal(msg)
    usm.logger.Printf("Sending to user %s: %s", userID, string(msgBytes))

    var lastErr error
    for _, conn := range conns {
        if err := conn.Conn.WriteJSON(msg); err != nil {
            usm.logger.Printf("Error sending to connection: %v", err)
            lastErr = err
        }
    }

    return lastErr
}

// Broadcast sends a message to all connected users
func (usm *UserSocketManager) Broadcast(data interface{}) {
    usm.mu.RLock()
    defer usm.mu.RUnlock()

    msg := map[string]interface{}{
        "type":      "broadcast",
        "timestamp": time.Now(),
        "data":      data,
    }

    // Log what we're broadcasting
    msgBytes, _ := json.Marshal(msg)
    usm.logger.Printf("Broadcasting: %s", string(msgBytes))

    for userID, conns := range usm.connections {
        for _, conn := range conns {
            if err := conn.Conn.WriteJSON(msg); err != nil {
                usm.logger.Printf("Error broadcasting to user %s: %v", userID, err)
            }
        }
    }
}

// Info returns information about active connections
func (usm *UserSocketManager) Info() map[string]interface{} {
    usm.mu.RLock()
    defer usm.mu.RUnlock()

    userCount := len(usm.connections)
    
    totalConnections := 0
    userConnections := make(map[string]int)
    
    for userID, conns := range usm.connections {
        userConnections[userID] = len(conns)
        totalConnections += len(conns)
    }

    return map[string]interface{}{
        "users":               userCount,
        "totalConnections":    totalConnections,
        "connectionsPerUser":  userConnections,
        "timestamp":           time.Now(),
    }
}

// Status logs and returns the current status of the websocket manager
func (usm *UserSocketManager) Status() map[string]interface{} {
    info := usm.Info()
    
    // Log the status
    infoBytes, _ := json.Marshal(info)
    usm.logger.Printf("WebSocket Status: %s", string(infoBytes))
    
    return info
}

// Global instance of UserSocketManager
var Manager *UserSocketManager

// Initialize the global UserSocketManager
func Init(logger *log.Logger) {
    Manager = NewUserSocketManager(logger)
}

// GetManager returns the global UserSocketManager instance
func GetManager() *UserSocketManager {
    return Manager
}