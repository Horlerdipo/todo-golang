package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/horlerdipo/todo-golang/env"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/users"
	"github.com/horlerdipo/todo-golang/utils"
	"log"
	"strconv"
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
	//check if email exists
	user, err := service.UserRepository.FindUserByEmail(email)
	if err != nil {
		return dtos.LoginUserResponseDto{}, errors.New("email or password is not valid")
	}

	//check if password is correct
	status := utils.CheckPasswordHash(password, user.Password)
	if status == false {
		return dtos.LoginUserResponseDto{}, errors.New("email or password is not valid")
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
