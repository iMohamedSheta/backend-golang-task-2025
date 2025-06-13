package services

import (
	"errors"
	"taskgo/internal/api/requests"
	"taskgo/internal/database/models"
	"taskgo/internal/enums"
	"taskgo/internal/repository"
	pkgErrors "taskgo/pkg/errors"
	"taskgo/pkg/utils"

	"gorm.io/gorm"
)

type UserService struct {
	userRepository *repository.UserRepository
}

// Create a new user service
func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

/*
 * Create a new user
 * @param req
 * @return *models.User, error[validationErr, error]
 */
func (s *UserService) CreateUser(req *requests.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Email:       req.Email,
		Password:    req.Password,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Role:        enums.RoleCustomer,
	}

	err := s.userRepository.Create(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Get user by id
func (s *UserService) GetUserById(id string) (*models.User, error) {
	return s.userRepository.FindById(id)
}

// Update user
func (s *UserService) UpdateUserAndGet(id string, req *requests.UpdateUserRequest) (*models.User, error) {
	updatedFields := make(map[string]any)
	sentFields := req.GetRequestSentFields()
	validKeys := utils.GetJSONKeys(req)

	for key, value := range sentFields {
		if validKeys[key] {
			updatedFields[key] = value
		}
	}

	user, err := s.userRepository.UpdateByIdAndGet(id, updatedFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewNotFoundError("user not found", "user not found", err)
		}
		return nil, err
	}

	return user, nil
}
