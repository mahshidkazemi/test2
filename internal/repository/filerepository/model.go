package filerepository

// // Queue represents a message queue
// type Queue struct {
// 	ID        string         // Unique identifier for the queue
// 	Name      string         // Name of the queue
// 	Values    []QueueData    // The actual queue data
// 	Consumers map[string]int // Map of consumer IDs to their read positions
// 	mu        sync.RWMutex   // Mutex for concurrent access
// }

// type QueueData struct {
// 	ID    string
// 	Index int64
// 	Value int64
// }

// // NewQueue creates a new queue
// func NewQueue(id, name string) *Queue {
// 	return &Queue{
// 		ID:        id,
// 		Name:      name,
// 		Values:    make([]QueueData, 0),
// 		Consumers: make(map[string]int),
// 	}
// }

// // Append adds a value to the queue
// func (q *Queue) Append(data QueueData) {
// 	q.mu.Lock()
// 	defer q.mu.Unlock()

// 	q.Values = append(q.Values, data)
// }

// // Read reads a value from the queue for a consumer
// // Returns the value and a boolean indicating if there was a value to read
// func (q *Queue) Read(consumerID string) (QueueData, bool) {
// 	q.mu.Lock()
// 	defer q.mu.Unlock()

// 	// If no values, return immediately
// 	if len(q.Values) == 0 {
// 		return QueueData{}, false
// 	}

// 	// Get the position for this consumer, or 0 if this is a new consumer
// 	position, exists := q.Consumers[consumerID]
// 	if !exists {
// 		position = 0
// 		q.Consumers[consumerID] = position
// 	}

// 	// If we've read all values, return false
// 	if position >= len(q.Values) {
// 		return QueueData{}, false
// 	}

// 	// Get the value at the current position
// 	value := q.Values[position]

// 	// Increment the position
// 	q.Consumers[consumerID]++

// 	return value, true
// }

// // ReadAt reads a value from the queue
// func (q *Queue) ReadAt(index int) (QueueData, error) {
// 	q.mu.Lock()
// 	defer q.mu.Unlock()

// 	// If no values, return immediately
// 	if len(q.Values) == 0 {
// 		return QueueData{}, repository.ErrQueueIsEmpty
// 	}

// 	// If we've read all values, return false
// 	if index >= len(q.Values) {
// 		return QueueData{}, repository.ErrIndexOutOfRange
// 	}

// 	// Get the value at the current position
// 	value := q.Values[index]

// 	return value, nil
// }

// // GetPosition returns the current read position for a consumer
// func (q *Queue) GetPosition(consumerID string) int {
// 	q.mu.RLock()
// 	defer q.mu.RUnlock()

// 	position, exists := q.Consumers[consumerID]
// 	if !exists {
// 		return -1
// 	}

// 	return position
// }

// // SetPosition sets the read position for a consumer
// func (q *Queue) SetPosition(consumerID string, position int) {
// 	q.mu.Lock()
// 	defer q.mu.Unlock()

// 	q.Consumers[consumerID] = position
// }

// func (q *Queue) AdvancePosition(consumerID string) int {
// 	q.mu.RLock()
// 	defer q.mu.RUnlock()

// 	_, exists := q.Consumers[consumerID]
// 	if !exists {
// 		q.Consumers[consumerID] = 0
// 		return 0
// 	}

// 	// Increment the position
// 	q.Consumers[consumerID]++

// 	return q.Consumers[consumerID]
// }

// // Size returns the number of values in the queue
// func (q *Queue) Size() int {
// 	q.mu.RLock()
// 	defer q.mu.RUnlock()

// 	return len(q.Values)
// }

// // =============================================================================

// type Clients struct {
// 	IDs []string
// 	mu  sync.RWMutex // Mutex for concurrent access
// }

// // NewClients creates a new Clients instance
// func NewClients() *Clients {
// 	return &Clients{
// 		IDs: make([]string, 0),
// 	}
// }

// func (c *Clients) Add(clientID string) error {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	for _, id := range c.IDs {
// 		if id == clientID {
// 			return fmt.Errorf("id exist") // ID exists
// 		}
// 	}

// 	c.IDs = append(c.IDs, clientID)

// 	return nil
// }

// // Remove removes a client ID from the list
// // Returns true if the ID was found and removed, false otherwise
// func (c *Clients) Remove(clientID string) bool {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	for i, id := range c.IDs {
// 		if id == clientID {
// 			c.IDs = append(c.IDs[:i], c.IDs[i+1:]...)
// 			return true
// 		}
// 	}
// 	return false
// }

// // Contains checks if a client ID exists in the list
// func (c *Clients) Contains(clientID string) bool {
// 	c.mu.RLock()
// 	defer c.mu.RUnlock()

// 	for _, id := range c.IDs {
// 		if id == clientID {
// 			return true
// 		}
// 	}
// 	return false
// }

// // GetAll returns a copy of all client IDs
// func (c *Clients) GetAll() []string {
// 	c.mu.RLock()
// 	defer c.mu.RUnlock()

// 	// Return a copy to prevent modification of the original slice
// 	ids := make([]string, len(c.IDs))
// 	copy(ids, c.IDs)
// 	return ids
// }

// // Size returns the number of client IDs
// func (c *Clients) Size() int {
// 	c.mu.RLock()
// 	defer c.mu.RUnlock()

// 	return len(c.IDs)
// }

// // Clear removes all client IDs
// func (c *Clients) Clear() {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	c.IDs = make([]string, 0)
// }
