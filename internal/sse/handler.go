package sse

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/middlewares"
	"net/http"
	"time"
)

type Handler struct {
	SSEService *Service
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares.JwtAuthMiddleware(h.SSEService.TokenBlacklistRepository))
		r.Get("/sse", h.registerSSE)
	})
}

func NewHandler(sseService *Service) *Handler {
	return &Handler{
		SSEService: sseService,
	}
}

func (h *Handler) registerSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rc := http.NewResponseController(w)
	_, err := fmt.Fprintf(w, "data: {\"type\":\"connected\"}\n\n")
	if err != nil {
		return
	}
	err = rc.Flush()
	if err != nil {
		return
	}

	clientGone := r.Context().Done()
	keepAlive := time.NewTicker(time.Second * 10)
	defer keepAlive.Stop()
	userId := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails).UserId

	client := &ConnectedClient{
		Data:   make(chan dtos.SSEData),
		Writer: w,
		Done:   clientGone,
		Quit:   make(chan struct{}),
	}
	fmt.Println("Registering SSE client", client)
	h.SSEService.AddClient(userId, client)

	for {
		select {
		case quitMsg := <-client.Quit:
			fmt.Println("Quit Message received", quitMsg)
			return
		case <-client.Done:
			fmt.Println("clientGone")
			return
		case <-keepAlive.C:
			_, err := fmt.Fprintf(w, ": heartbeat\n\n")
			if err != nil {
				fmt.Println(err)
				return
			}
			err = rc.Flush()
			if err != nil {
				fmt.Println(err)
				return
			}
		case msg := <-client.Data:
			fmt.Println("sending message on handler", msg)
			marshalledMsg, err := json.Marshal(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = fmt.Fprintf(w, "data: %s\n\n", marshalledMsg)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = rc.Flush()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
