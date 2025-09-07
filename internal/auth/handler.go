package auth

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/utils"
	"log"
	"net/http"
)

type Handler struct {
	AuthService *Service
}

func NewAuthHandler(authService *Service) *Handler {
	return &Handler{
		AuthService: authService,
	}
}

func (h *Handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Login Successful"}`))
}

func (h *Handler) registerHandler(w http.ResponseWriter, r *http.Request) {

	createUserDto, err := utils.JsonValidate[dtos.CreateUserDTO](w, r)
	if err != nil {
		return
	}

	//send to auth service
	_, err = h.AuthService.Register(createUserDto)
	if err != nil {
		log.Println(fmt.Sprintf("Error while creating user: %v", err))
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), nil)
	}

	//return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.loginHandler)
		r.Post("/register", h.registerHandler)
	})
}
