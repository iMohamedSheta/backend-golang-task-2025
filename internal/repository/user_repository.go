package repository

import (
	"errors"
	"taskgo/internal/database/models"
	"taskgo/internal/deps"
	"time"
)

type UserRepository struct {
	db *deps.GormDB
}

func NewUserRepository(db *deps.GormDB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update last login time for a user
func (r *UserRepository) UpdateLastLogin(user *models.User) error {
	user.LastLoginAt = time.Now()
	return r.db.DB.Save(user).Error
}

// Create a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.DB.Create(user).Error
}

// Get a user by id
func (r *UserRepository) FindById(id string) (*models.User, error) {
	var user models.User
	if id == "" {
		return nil, errors.New("id is required")
	}

	if err := r.db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update user
func (r *UserRepository) UpdateById(id string, data map[string]interface{}) error {
	if id == "" {
		return errors.New("id is required")
	}

	if err := r.db.DB.Model(&models.User{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// Update user and return the updated user
func (r *UserRepository) UpdateByIdAndGet(id string, data map[string]interface{}) (*models.User, error) {
	err := r.UpdateById(id, data)

	if err != nil {
		return nil, err
	}

	// Fetch the updated user to return the latest state
	var updatedUser models.User
	if err := r.db.DB.Where("id = ?", id).First(&updatedUser).Error; err != nil {
		return nil, err
	}

	return &updatedUser, nil
}
