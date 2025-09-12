package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/horlerdipo/todo-golang/env"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/users"
	"github.com/horlerdipo/todo-golang/pkg"
	"github.com/horlerdipo/todo-golang/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	UserRepository users.UserRepository
}

func NewService(userRepository users.UserRepository) *Service {
	return &Service{
		UserRepository: userRepository,
	}
}

func (service *Service) Register(userDto dtos.CreateUserDTO) (bool, error) {
	//check if user already exists
	userDto.Email = strings.ToLower(userDto.Email)

	_, err := service.UserRepository.FindUserByEmail(userDto.Email)
	if err == nil {
		return false, errors.New("user already exists")
	}

	// create user
	_, err = service.UserRepository.CreateUser(&userDto)
	if err != nil {
		log.Println("Error while creating user: " + err.Error())
		return false, err
	}

	//send success
	return true, nil
}

func (service *Service) Login(email string, password string) (dtos.LoginUserResponseDto, error) {
	email = strings.ToLower(email)

	//check if email exists
	user, err := service.UserRepository.FindUserByEmail(email)
	if err != nil {
		return dtos.LoginUserResponseDto{}, errors.New("email or password is not valid")
	}

	//check if password is correct
	status := utils.CheckPasswordHash(password, user.Password)
	if status == false {
		return dtos.LoginUserResponseDto{}, errors.New("email or password is not correct")
	}

	//generate jwt token
	ttlEnv := env.FetchInt("JWT_TTL", 24)
	ttl := time.Now().Add(time.Duration(ttlEnv) * time.Hour)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     ttl.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString := ""
	tokenString, err = token.SignedString([]byte(env.FetchString("JWT_SECRET")))

	if err != nil {
		return dtos.LoginUserResponseDto{}, errors.New("error while signing token")
	}

	return dtos.LoginUserResponseDto{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Token: dtos.TokenDetails{
			Token: tokenString,
			Exp:   strconv.FormatInt(ttl.Unix(), 10),
		},
	}, nil
}

func (service *Service) SendForgotPasswordToken(email string) (bool, error) {
	//check if email exists
	email = strings.ToLower(email)

	user, err := service.UserRepository.FindUserByEmail(email)
	if err != nil {
		return false, errors.New("email does not exist")
	}

	//generate and save token
	resetToken := ""
	resetToken, err = utils.RandomNumericString(env.FetchInt("PASSWORD_RESET_TOKEN_LENGTH", 6))
	if err != nil {
		return false, errors.New("error while generating reset token")
	}

	resetTokenExpiresAt := time.Now().Add(time.Duration(env.FetchInt("PASSWORD_RESET_TOKEN_TTL")) * time.Minute)
	err = service.UserRepository.UpdateUser(user.ID, &dtos.UpdateUserDTO{
		ResetToken:          &resetToken,
		ResetTokenExpiresAt: &resetTokenExpiresAt,
	})

	if err != nil {
		return false, errors.New("error while updating reset token")
	}

	//send token via email
	go func() {
		err := pkg.SendEmail(pkg.SendEmailConfig{
			Recipients:  []string{email},
			Subject:     "Password Reset",
			Content:     "Hello, your password reset token is " + resetToken + " and it expires by " + resetTokenExpiresAt.Format("2006-01-02 15:04:05"),
			ContentType: "text/plain",
		})
		if err != nil {
			log.Println(err)
		}
	}()
	//return success
	return true, nil
}

func (service *Service) ResetPassword(resetToken string, newPassword string) error {
	//check if reset password token exists
	user, err := service.UserRepository.FindUserByResetToken(resetToken)
	if err != nil {
		return errors.New("reset token is invalid")
	}

	//check if it has not expired
	if user.ResetTokenExpiresAt == nil {
		return errors.New("reset token has expired")
	}
	if time.Now().After(*user.ResetTokenExpiresAt) {
		return errors.New("reset token has expired")
	}

	//update the password
	err = service.UserRepository.UpdateUserPassword(user.ID, newPassword, true)
	if err != nil {
		log.Println("Error while updating user password: " + err.Error())
		return errors.New("error while resetting password")
	}

	return nil
}
