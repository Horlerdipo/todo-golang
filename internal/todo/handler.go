package todo

import (
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/middlewares"
	"github.com/horlerdipo/todo-golang/utils"
	"net/http"
)

type Handler struct {
	TodoService *Service
}

func NewHandler(todoService *Service) *Handler {
	return &Handler{
		TodoService: todoService,
	}
}

func (handler *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/todos", func(r chi.Router) {
		r.Use(middlewares.JwtAuthMiddleware(handler.TodoService.TokenBlacklistRepository))
		r.Post("/", handler.CreateTodo)
	})
}

func (handler *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	jsonResponse, err := utils.JsonValidate[dtos.CreateTodoDTO](w, r)
	if err != nil {
		return
	}

	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)
	jsonResponse.UserID = authDetails.UserId

	_, err = handler.TodoService.CreateTodo(&jsonResponse)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
