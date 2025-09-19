package database

import (
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/utils"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

type UserRepository interface {
	FindUserByID(ctx context.Context, id uint) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, userDto *dtos.CreateUserDTO) (uint, error)
	UpdateUser(ctx context.Context, userId uint, userDto *dtos.UpdateUserDTO) error
	FindUserByResetToken(ctx context.Context, token string) (*User, error)
	UpdateUserPassword(ctx context.Context, userId uint, password string, resetTokens bool) error
}

type userRepository struct {
	db *gorm.DB
}

func (repo *userRepository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	userModel := User{}
	result := repo.db.WithContext(ctx).Where("email = ?", email).First(&userModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userModel, nil
}

func (repo *userRepository) CreateUser(ctx context.Context, userDto *dtos.CreateUserDTO) (uint, error) {
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

	result := repo.db.WithContext(ctx).Create(&userModel)

	if result.Error != nil {
		return 0, result.Error
	}
	return userModel.ID, nil
}

func (repo *userRepository) UpdateUser(ctx context.Context, userId uint, userDto *dtos.UpdateUserDTO) error {
	result := repo.db.WithContext(ctx).Model(&User{}).Where("id = ?", userId).Updates(userDto)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *userRepository) FindUserByResetToken(ctx context.Context, resetToken string) (*User, error) {
	userModel := User{}
	result := repo.db.WithContext(ctx).Where("reset_token = ?", resetToken).First(&userModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userModel, nil
}

func (repo *userRepository) UpdateUserPassword(ctx context.Context, userId uint, password string, resetTokens bool) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	if resetTokens {
		result := repo.db.WithContext(ctx).Model(&User{}).Where("id = ?", userId).Updates(map[string]interface{}{"password": hashedPassword, "reset_token": nil, "reset_token_expires_at": nil})
		if result.Error != nil {
			return result.Error
		}
	} else {
		result := repo.db.WithContext(ctx).Model(&User{}).Where("id = ?", userId).Update("password", hashedPassword)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (repo *userRepository) FindUserByID(ctx context.Context, id uint) (*User, error) {
	userModel := User{}
	result := repo.db.WithContext(ctx).Where("id = ?", id).First(&userModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userModel, nil
}
