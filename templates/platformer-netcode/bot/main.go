// Package main implements a bot client for load testing the drip server.
package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	// Import types for actions and messages.
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/actions"
)

// Constants for connection and behavior parameters.
const (
	BOT_COUNT = 300
	ADDR      = "localhost"

	// connectionTimeout defines time limit for establishing a connection.
	connectionTimeout = 5 * time.Second
	// readDeadline defines time limit for read operations.
	readDeadline = 10 * time.Second
	// writeDeadline defines time limit for write operations.
	writeDeadline = 2 * time.Second
	// monitorCheckInterval defines frequency of bot status checks.
	monitorCheckInterval = 500 * time.Millisecond
	// shutdownWaitTimeout defines maximum wait time for bot termination.
	shutdownWaitTimeout = 10 * time.Second
	// lengthPrefixBytes defines size of message length prefix.
	lengthPrefixBytes = 4

	// actionSendInterval controls action send frequency.
	actionSendInterval = 100 * time.Millisecond
	// minStateDuration defines minimum time in movement/idle state.
	minStateDuration = 500 * time.Millisecond
	// maxStateDuration defines maximum time in movement/idle state.
	maxStateDuration = 3000 * time.Millisecond
	// downChance defines probability of sending Down action per interval.
	downChance = 0.1
)

// BotState represents current high-level behavior.
type BotState int

const (
	StateIdle BotState = iota
	StateMovingLeft
	StateMovingRight
)

// BotClient manages state and network connection for a bot instance.
type BotClient struct {
	id           int
	conn         net.Conn
	currentState BotState
	stateEndTime time.Time // Time when state should change.

	running     bool
	mutex       sync.Mutex // Protects conn, running, state vars, actionStamp.
	actionStamp int        // Monotonically increasing counter for sent actions.
}

// NewBotClient creates and initializes a connected bot client.
func NewBotClient(id int, serverAddr string) (*BotClient, error) {
	log.Printf("[Bot %d] Connecting to %s", id, serverAddr)
	conn, err := net.DialTimeout("tcp", serverAddr, connectionTimeout)
	if err != nil {
		log.Printf("[Bot %d] Failed connection: %v", id, err)
		return nil, fmt.Errorf("bot %d connection failed", id)
	}
	log.Printf("[Bot %d] Connected to %s", id, conn.RemoteAddr())

	// Initialize state.
	return &BotClient{
		id:           id,
		conn:         conn,
		currentState: StateIdle,
		stateEndTime: time.Now(),
		running:      true,
		actionStamp:  0,
	}, nil
}

// Start launches bot's action and monitoring goroutines.
func (b *BotClient) Start() {
	b.mutex.Lock()
	if !b.running {
		b.mutex.Unlock()
		log.Printf("[Bot %d] Start called but already stopped.", b.id)
		return
	}
	b.mutex.Unlock()

	log.Printf("[Bot %d] Starting loops", b.id)
	go b.actionLoop()
	go b.connectionMonitor()
}

// Stop shuts down the bot and closes its connection.
func (b *BotClient) Stop() {
	b.mutex.Lock()
	if !b.running {
		b.mutex.Unlock()
		return
	}
	b.running = false
	log.Printf("[Bot %d] Stopping...", b.id)
	connToClose := b.conn
	b.conn = nil
	b.mutex.Unlock()

	if connToClose != nil {
		err := connToClose.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("[Bot %d] Error closing connection: %v", b.id, err)
		}
	}
}

// IsRunning reports whether the bot is currently marked as running.
func (b *BotClient) IsRunning() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.running
}

// connectionMonitor checks connection health and stops the bot on errors.
func (b *BotClient) connectionMonitor() {
	log.Printf("[Bot %d] Connection monitor started.", b.id)
	buffer := make([]byte, 1)
	defer log.Printf("[Bot %d] Connection monitor finished.", b.id)

	for {
		if !b.IsRunning() {
			return
		}

		b.mutex.Lock()
		currentConn := b.conn
		b.mutex.Unlock()
		if currentConn == nil {
			return
		}

		err := currentConn.SetReadDeadline(time.Now().Add(readDeadline))
		if err != nil {
			if b.IsRunning() {
				log.Printf("[Bot %d] Monitor set deadline error: %v. Stopping.", b.id, err)
				b.Stop()
			}
			return
		}

		_, err = currentConn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if b.IsRunning() {
				log.Printf("[Bot %d] Monitor read error: %v. Stopping.", b.id, err)
				b.Stop()
			}
			return
		}
	}
}

// chooseNewState selects next state randomly and sets its duration.
func (b *BotClient) chooseNewState() {
	rnd := rand.Float32()
	newState := StateIdle

	if rnd < 0.35 {
		newState = StateMovingLeft
	} else if rnd < 0.70 {
		newState = StateMovingRight
	}

	durationRange := int64(maxStateDuration - minStateDuration)
	duration := minStateDuration + time.Duration(rand.Int63n(durationRange))
	endTime := time.Now().Add(duration)

	b.mutex.Lock()
	b.currentState = newState
	b.stateEndTime = endTime
	b.mutex.Unlock()

	log.Printf("[Bot %d] Entering state: %v for %v", b.id, newState, duration)
}

// actionLoop manages state transitions and sends actions based on current state.
func (b *BotClient) actionLoop() {
	log.Printf("[Bot %d] Action loop started.", b.id)
	ticker := time.NewTicker(actionSendInterval)
	defer ticker.Stop()
	defer log.Printf("[Bot %d] Action loop finished.", b.id)

	for {
		select {
		case <-ticker.C:
		case <-time.After(monitorCheckInterval):
			if !b.IsRunning() {
				return
			}
			continue
		}

		if !b.IsRunning() {
			return
		}

		now := time.Now()

		b.mutex.Lock()
		endTime := b.stateEndTime
		currentState := b.currentState
		b.mutex.Unlock()

		if now.After(endTime) {
			b.chooseNewState()
			b.mutex.Lock()
			currentState = b.currentState
			b.mutex.Unlock()
		}

		shouldSendAction := false
		var actionToSend input.Action

		switch currentState {
		case StateMovingLeft:
			actionToSend = actions.Left
			shouldSendAction = true
		case StateMovingRight:
			actionToSend = actions.Right
			shouldSendAction = true
		case StateIdle:
			shouldSendAction = false
		}

		rnd := rand.Float32()
		if rnd < downChance {
			actionToSend = actions.Down
			shouldSendAction = true
		}

		if shouldSendAction {
			b.mutex.Lock()
			currentStamp := b.actionStamp
			b.actionStamp++
			currentConn := b.conn
			isRunning := b.running
			b.mutex.Unlock()

			if !isRunning {
				return
			}
			if currentConn == nil {
				continue
			}

			stampedAction := input.StampedAction{Val: actionToSend, Tick: currentStamp}

			actionMsg := input.ClientActionMessage{ReceiverIndex: 0, Actions: []input.StampedAction{stampedAction}}
			msgData, err := json.Marshal(actionMsg)
			if err != nil {
				log.Printf("[Bot %d] Marshal error for action %v: %v. Skipping.", b.id, actionToSend, err)
				continue
			}

			msgLen := uint32(len(msgData))
			prefixBytes := make([]byte, lengthPrefixBytes)
			binary.BigEndian.PutUint32(prefixBytes, msgLen)
			framedMsg := append(prefixBytes, msgData...)

			err = currentConn.SetWriteDeadline(time.Now().Add(writeDeadline))
			if err != nil {
				if b.IsRunning() {
					log.Printf("[Bot %d] Set deadline error: %v. Stopping.", b.id, err)
					b.Stop()
				}
				return
			}

			_, err = currentConn.Write(framedMsg)
			if err != nil {
				if b.IsRunning() {
					log.Printf("[Bot %d] Send error for action %v: %v. Stopping.", b.id, actionToSend, err)
					b.Stop()
				}
				return
			}
		}
	}
}

func main() {
	// Parse command line flags.
	numBots := flag.Int("bots", BOT_COUNT, "Number of bot clients to create")
	serverAddr := flag.String("server", ADDR+":8080", "Server address (host:port)")
	flag.Parse()

	log.Printf("--- Bot Swarm Starting ---")
	log.Printf("Server: %s, Bots: %d", *serverAddr, *numBots)

	rand.New(rand.NewSource(time.Now().UnixNano()))
	bots := make([]*BotClient, 0, *numBots)
	var wg sync.WaitGroup

	log.Printf("Launching bots...")
	launchedCount := 0
	for i := 0; i < *numBots; i++ {
		bot, err := NewBotClient(i, *serverAddr)
		if err != nil {
			continue
		}
		bots = append(bots, bot)
		launchedCount++
		wg.Add(1)
		go func(b *BotClient) {
			defer wg.Done()
			b.Start()
			monitorTicker := time.NewTicker(monitorCheckInterval)
			defer monitorTicker.Stop()
			for range monitorTicker.C {
				if !b.IsRunning() {
					break
				}
			}
		}(bot)
		time.Sleep(time.Duration(rand.Intn(50)+50) * time.Millisecond)
	}
	log.Printf("--- %d bots launched ---", launchedCount)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Bot swarm running. Press Ctrl+C to stop.")

	sig := <-sigChan
	log.Printf("--- Received signal: %v. Shutting down... ---", sig)

	log.Printf("Stopping %d bots...", len(bots))
	for _, bot := range bots {
		bot.Stop()
	}

	log.Println("Waiting for bots to finish...")
	waitChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitChan)
	}()

	select {
	case <-waitChan:
		log.Println("--- All bot goroutines finished gracefully. ---")
	case <-time.After(shutdownWaitTimeout):
		log.Println("--- WARNING: Timeout waiting for bots to finish. ---")
	}
	log.Println("--- Bot swarm shutdown complete ---")
}

// String returns a human-readable representation of the BotState.
func (s BotState) String() string {
	switch s {
	case StateIdle:
		return "Idle"
	case StateMovingLeft:
		return "MovingLeft"
	case StateMovingRight:
		return "MovingRight"
	default:
		return fmt.Sprintf("UnknownState(%d)", int(s))
	}
}
