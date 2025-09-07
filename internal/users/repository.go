package users

import (
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/utils"
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

type UserRepository interface {
	FindUserByEmail(email string) (*User, error)
	CreateUser(userDto *dtos.CreateUserDTO) (uint, error)
}

type userRepository struct {
	db *gorm.DB
}

func (repo *userRepository) FindUserByEmail(email string) (*User, error) {
	userModel := User{}
	result := repo.db.Where("email = ?", email).First(&userModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userModel, nil
}

func (repo *userRepository) CreateUser(userDto *dtos.CreateUserDTO) (uint, error) {
	hashedPassword, err := utils.HashPassword(userDto.Password)
	if err != nil {
		return 0, err
	}

	userModel := User{
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
		Email:     userDto.Email,
		Password:  hashedPassword,
	}

	result := repo.db.Create(&userModel)

	if result.Error != nil {
		return 0, result.Error
	}
	return userModel.ID, nil
}
