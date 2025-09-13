package auth

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/middlewares"
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
	loginDto, err := utils.JsonValidate[dtos.LoginUserDTO](w, r)
	if err != nil {
		log.Println(err)
		return
	}

	response, err := h.AuthService.Login(loginDto.Email, loginDto.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, "login successful", response)
	return
}

func (h *Handler) registerHandler(w http.ResponseWriter, r *http.Request) {

	createUserDto, err := utils.JsonValidate[dtos.CreateUserDTO](w, r)
	if err != nil {
		log.Println(err)
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

func (h *Handler) sendResetPasswordToken(w http.ResponseWriter, r *http.Request) {
	//validate
	type EmailStruct struct {
		Email string `json:"email" validate:"required"`
	}

	email, err := utils.JsonValidate[EmailStruct](w, r)
	if err != nil {
		return
	}

	//call service
	resp, errorRsp := h.AuthService.SendForgotPasswordToken(email.Email)
	if errorRsp != nil && resp == false {
		utils.RespondWithError(w, http.StatusBadRequest, errorRsp.Error(), nil)
		return
	}

	//return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Handler) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	type ResetPassword struct {
		NewPassword string `json:"new_password" validate:"required"`
		ResetToken  string `json:"reset_token" validate:"required"`
	}

	resetPasswordStruct, err := utils.JsonValidate[ResetPassword](w, r)
	if err != nil {
		return
	}

	err = h.AuthService.ResetPassword(resetPasswordStruct.ResetToken, resetPasswordStruct.NewPassword)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Handler) profileHandler(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middlewares.UserKey)
	if userIDValue == nil {
		utils.RespondWithError(w, 400, "user not found", nil)
		return
	}

	userId, ok := userIDValue.(uint)
	if !ok {
		utils.RespondWithError(w, 400, "invalid user type", nil)
		return
	}

	user, err := h.AuthService.FetchUserDetails(userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, "profile fetched", user)
	return
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.loginHandler)
		r.Post("/register", h.registerHandler)
		r.Post("/password/forgot", h.sendResetPasswordToken)
		r.Post("/password/reset", h.resetPasswordHandler)
		r.With(middlewares.JwtAuthMiddleware).Get("/user", h.profileHandler)
	})
}
