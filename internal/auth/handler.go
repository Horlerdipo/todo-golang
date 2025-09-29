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

	response, err := h.AuthService.Login(r.Context(), loginDto.Email, loginDto.Password)
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
	_, err = h.AuthService.Register(r.Context(), createUserDto)
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
	resp, errorRsp := h.AuthService.SendForgotPasswordToken(r.Context(), email.Email)
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

	err = h.AuthService.ResetPassword(r.Context(), resetPasswordStruct.ResetToken, resetPasswordStruct.NewPassword)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Handler) profileHandler(w http.ResponseWriter, r *http.Request) {
	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)

	user, err := h.AuthService.FetchUserDetails(r.Context(), authDetails.UserId)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, "profile fetched", user)
	return
}

func (h *Handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)
	resp := h.AuthService.LogoutUser(r.Context(), authDetails.UserId, authDetails.JwtToken, authDetails.JwtExpirationTime.Time)
	if !resp {
		utils.RespondWithError(w, http.StatusInternalServerError, "unable to log out, please try again", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.loginHandler)
		r.Post("/register", h.registerHandler)
		r.Post("/password/forgot", h.sendResetPasswordToken)
		r.Post("/password/reset", h.resetPasswordHandler)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.JwtAuthMiddleware(h.AuthService.TokenBlacklistRepository))
			r.Get("/user", h.profileHandler)
			r.Post("/logout", h.logoutHandler)
		})
	})
}
