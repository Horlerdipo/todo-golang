package sse

import (
	"fmt"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"net/http"
	"sync"
)

type ConnectedClient struct {
	Data   chan dtos.SSEData
	Writer http.ResponseWriter
	Quit   chan struct{}
	Done   <-chan struct{}
}

type Service struct {
	mutex                    sync.RWMutex
	ConnectedClients         map[uint][]*ConnectedClient
	TokenBlacklistRepository database.TokenBlacklistRepository
}

func NewService(tokenBlacklistRepository database.TokenBlacklistRepository) *Service {
	connectedClient := make(map[uint][]*ConnectedClient)
	return &Service{
		ConnectedClients:         connectedClient,
		TokenBlacklistRepository: tokenBlacklistRepository,
	}
}

func (service *Service) AddClient(userId uint, client *ConnectedClient) {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	fmt.Println("Adding client", userId, client)
	service.ConnectedClients[userId] = append(service.ConnectedClients[userId], client)
	fmt.Printf("Client added %v", service.ConnectedClients)
}

func (service *Service) RemoveClients(userId uint) {
	fmt.Printf("Client before removal %v", service.ConnectedClients)
	service.mutex.RLock()
	clients := service.ConnectedClients[userId]
	service.mutex.RUnlock()

	for _, client := range clients {
		close(client.Quit)
	}

	if _, ok := service.ConnectedClients[userId]; ok {
		service.mutex.RLock()
		delete(service.ConnectedClients, userId)
		service.mutex.RUnlock()
	}
	fmt.Printf("Client removed %v", service.ConnectedClients)

}

func (service *Service) SendMessage(userId uint, message dtos.SSEData) {
	service.mutex.RLock()
	clients := service.ConnectedClients[userId]
	service.mutex.RUnlock()

	fmt.Println("\nabout to send message to", len(clients), "clients")
	for _, client := range clients {
		fmt.Println("sending message", client)
		select {
		case client.Data <- message:
			// Message sent successfully
			fmt.Println("Message sent")
		case <-client.Done:
			// Client disconnected, skip
			fmt.Println("Client already disconnected")
		default:
			// Channel full or blocked - don't block, just skip
			fmt.Println("Channel blocked or full, skipping")
		}
	}
}
