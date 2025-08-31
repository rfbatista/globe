package globe

import (
	"errors"
	"sync"
)

type Mailbox interface {
	PostSystemMessage(Message)
	PostUserMessage(Message)
	RegisterActor(*worker)
	RegisterInboundChannel(chan Message)
	RegisterActorChannel(chan Message)
	Start() error
}

type DefaultMailbox struct {
	systemMessages chan Message
	userMessages   chan Message

	actor          *worker
	inboundChannel chan Message
	actorChannel   chan Message

	mu      sync.RWMutex
	started bool
	stopCh  chan struct{}
	wg      sync.WaitGroup
}

func NewMailbox() Mailbox {
	return &DefaultMailbox{
		systemMessages: make(chan Message, 128),
		userMessages:   make(chan Message, 128),
		stopCh:         make(chan struct{}),
	}
}

// Post a system-level message
func (m *DefaultMailbox) PostSystemMessage(msg Message) {
	m.systemMessages <- msg
}

// Post a user-level message
func (m *DefaultMailbox) PostUserMessage(msg Message) {
	m.userMessages <- msg
}

// Register the actor (worker) that consumes messages
func (m *DefaultMailbox) RegisterActor(w *worker) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actor = w
}

// Register a channel to receive inbound messages (from outside world)
func (m *DefaultMailbox) RegisterInboundChannel(ch chan Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inboundChannel = ch
}

// Register the actor’s mailbox channel
func (m *DefaultMailbox) RegisterActorChannel(ch chan Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actorChannel = ch
}

// Start launches the dispatcher loop
func (m *DefaultMailbox) Start() error {
	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return errors.New("mailbox already started")
	}
	m.started = true
	m.mu.Unlock()

	m.wg.Add(1)
	go m.dispatch()

	return nil
}

func (m *DefaultMailbox) dispatch() {
	defer m.wg.Done()

	for {
		select {
		// system messages have priority
		case msg := <-m.systemMessages:
			m.deliver(msg)

		// user messages next
		case msg := <-m.userMessages:
			m.deliver(msg)

		// external inbound channel
		case msg := <-m.inboundChannel:
			m.deliver(msg)

		case <-m.stopCh:
			return
		}
	}
}

func (m *DefaultMailbox) deliver(msg Message) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.actorChannel != nil {
		m.actorChannel <- msg
	}
}

// Stop shuts down the mailbox dispatcher
func (m *DefaultMailbox) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.started {
		return
	}
	close(m.stopCh)
	m.wg.Wait()
	m.started = false
}
