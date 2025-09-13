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
	FindUserByID(id uint) (*User, error)
	FindUserByEmail(email string) (*User, error)
	CreateUser(userDto *dtos.CreateUserDTO) (uint, error)
	UpdateUser(userId uint, userDto *dtos.UpdateUserDTO) error
	FindUserByResetToken(token string) (*User, error)
	UpdateUserPassword(userId uint, password string, resetTokens bool) error
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

func (repo *userRepository) UpdateUser(userId uint, userDto *dtos.UpdateUserDTO) error {
	result := repo.db.Model(&User{}).Where("id = ?", userId).Updates(userDto)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *userRepository) FindUserByResetToken(resetToken string) (*User, error) {
	userModel := User{}
	result := repo.db.Where("reset_token = ?", resetToken).First(&userModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userModel, nil
}

func (repo *userRepository) UpdateUserPassword(userId uint, password string, resetTokens bool) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	if resetTokens {
		result := repo.db.Model(&User{}).Where("id = ?", userId).Updates(map[string]interface{}{"password": hashedPassword, "reset_token": nil, "reset_token_expires_at": nil})
		if result.Error != nil {
			return result.Error
		}
	} else {
		result := repo.db.Model(&User{}).Where("id = ?", userId).Update("password", hashedPassword)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (repo *userRepository) FindUserByID(id uint) (*User, error) {
	userModel := User{}
	result := repo.db.Where("id = ?", id).First(&userModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userModel, nil
}
