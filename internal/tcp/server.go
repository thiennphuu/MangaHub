package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"mangahub/pkg/models"
	"mangahub/pkg/utils"
)

// Server represents the TCP sync server
type Server struct {
	Port        string
	Connections map[string]net.Conn
	Broadcast   chan models.ProgressUpdate
	Register    chan net.Conn
	Unregister  chan net.Conn
	mutex       sync.RWMutex
	done        chan bool
	logger      *utils.Logger
}

// NewServer creates a new TCP server
func NewServer(port string, logger *utils.Logger) *Server {
	return &Server{
		Port:        port,
		Connections: make(map[string]net.Conn),
		Broadcast:   make(chan models.ProgressUpdate, 100),
		Register:    make(chan net.Conn),
		Unregister:  make(chan net.Conn),
		done:        make(chan bool),
		logger:      logger,
	}
}

// Start starts the TCP server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %w", err)
	}
	defer listener.Close()

	s.logger.Info("TCP server started on port %s", s.Port)

	// Start broadcast handler
	go s.handleBroadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("Error accepting connection: %v", err)
			continue
		}

		s.mutex.Lock()
		connID := fmt.Sprintf("conn_%d", len(s.Connections))
		s.Connections[connID] = conn
		s.mutex.Unlock()

		s.logger.Info("✅ New Connection | %-10s | %s", connID, conn.RemoteAddr())
		go s.handleConnection(connID, conn)
	}
}

// handleConnection handles a client connection
func (s *Server) handleConnection(connID string, conn net.Conn) {
	defer func() {
		s.mutex.Lock()
		delete(s.Connections, connID)
		s.mutex.Unlock()
		addr := conn.RemoteAddr()
		conn.Close()
		s.logger.Info("❌ Connection Closed | %-10s | %s", connID, addr)
	}()

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		var update models.ProgressUpdate
		if err := json.Unmarshal([]byte(line), &update); err != nil {
			s.logger.Error("Error parsing message: %v", err)
			continue
		}

		s.Broadcast <- update
	}
}

// handleBroadcast broadcasts progress updates to all clients
func (s *Server) handleBroadcast() {
	for update := range s.Broadcast {
		data, err := json.Marshal(update)
		if err != nil {
			s.logger.Error("Error marshaling update: %v", err)
			continue
		}

		s.mutex.RLock()
		for _, conn := range s.Connections {
			fmt.Fprintf(conn, "%s\n", data)
		}
		s.mutex.RUnlock()

		s.logger.Info("📡 Broadcast | User: %s | Manga: %s | Ch: %d", update.UserID, update.MangaID, update.Chapter)
	}
}

// Stop stops the server
func (s *Server) Stop() {
	s.mutex.Lock()
	for _, conn := range s.Connections {
		conn.Close()
	}
	s.Connections = make(map[string]net.Conn)
	s.mutex.Unlock()
	close(s.done)
}

// GetConnectionCount returns the number of active connections
func (s *Server) GetConnectionCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.Connections)
}
