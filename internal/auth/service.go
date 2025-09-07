package auth

import (
	"errors"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/users"
	"log"
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
