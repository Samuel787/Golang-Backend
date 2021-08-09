package middleware

import (
	"testing"

	"../models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mock *MockRepository) AddUser(user *models.User) (*models.User, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*models.User), args.Error(1)
}

func TestCreateUserPositive(t *testing.T) {
	mockRepo := new(MockRepository)
	user := models.User{
		Name:        "Jackson",
		DOB:         "21 Aug 2001",
		Address:     "Boon Lay",
		Description: "this is some dummy description",
		Latitude:    0.1,
		Longitude:   0.1}
	mockRepo.On("AddUser").Return(&user, nil)
	testService := NewUserService(mockRepo)
	err := testService.CreateUser(&user)
	mockRepo.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestCreateUserNoName(t *testing.T) {
	mockRepo := new(MockRepository)
	user := models.User{
		DOB:         "21 Aug 2001",
		Address:     "Boon Lay",
		Description: "this is some dummy description",
		Latitude:    0.1,
		Longitude:   0.1}
	mockRepo.On("AddUser").Return(&user, nil)
	testService := NewUserService(mockRepo)
	err := testService.CreateUser(&user)
	assert.NotNil(t, err)
	assert.Equal(t, "[Create User] Unable to create user with nil Name", err.Error())
}

func TestCreateUserNoLocation(t *testing.T) {
	mockRepo := new(MockRepository)
	user := models.User{
		Name:        "Jackson",
		DOB:         "21 Aug 2001",
		Address:     "Boon Lay",
		Description: "this is some dummy description"}
	mockRepo.On("AddUser").Return(&user, nil)
	testService := NewUserService(mockRepo)
	err := testService.CreateUser(&user)
	mockRepo.AssertExpectations(t)
	assert.Nil(t, err)
}
