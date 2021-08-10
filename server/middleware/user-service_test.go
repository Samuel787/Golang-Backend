package middleware

import (
	"errors"
	"testing"

	"../models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockRepository struct {
	mock.Mock
}

func (mock *MockRepository) AddUser(user *models.User) (*models.User, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*models.User), args.Error(1)
}

func (mock *MockRepository) GetAllUsers() ([]primitive.M, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]primitive.M), args.Error(1)
}

func (mock *MockRepository) DeleteUser(userId string) error {
	args := mock.Called()
	return args.Error(1)
}

func (mock *MockRepository) GetUser(user string) (bson.M, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(bson.M), args.Error(1)
}

func (mock *MockRepository) UpdateUser(userId string, update bson.M) error {
	args := mock.Called()
	return args.Error(1)
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

/**
Testing: GetAllUsers
Result: Should not throw any error even though there are no users in database
*/
func TestGetAllUsersSuccessNoUsers(t *testing.T) {
	mockRepo := new(MockRepository)
	var users []primitive.M
	mockRepo.On("GetAllUsers").Return(users, nil)
	testService := NewUserService(mockRepo)
	_, err := testService.GetAllUsers()
	assert.Nil(t, err)
}

/**
Testing: GetAllUsers
Result: Success, should not throw error when there are users in database
*/
func TestGetAllUsersSuccessWithUsers(t *testing.T) {
	mockRepo := new(MockRepository)
	var users []primitive.M
	users = append(users, primitive.M{
		"_id":         "61098bdd425f4ec901e268df",
		"address":     "Singapore",
		"createdAt":   "1100",
		"description": "I'm young",
		"dob":         "1 Jan 2020",
		"latitude":    9.9,
		"longitude":   11.11,
		"name":        "bean"})
	mockRepo.On("GetAllUsers").Return(users, nil)
	testService := NewUserService(mockRepo)
	usersResult, err := testService.GetAllUsers()
	assert.Nil(t, err)
	assert.Equal(t, users, usersResult)
}

/**
Testing: GetAllUsers
Result: Failure, database throws error for some unknown reason
*/
func TestGetAllUsersFailureDBError(t *testing.T) {
	mockRepo := new(MockRepository)
	var users []primitive.M
	mockError := errors.New("Dummy Error")
	mockRepo.On("GetAllUsers").Return(users, errors.New("Dummy Error"))
	testService := NewUserService(mockRepo)
	usersResult, err := testService.GetAllUsers()
	assert.Nil(t, usersResult)
	assert.NotNil(t, err)
	assert.Equal(t, mockError.Error(), err.Error())
}

/**
Testing: GetUserById
Result: Successful as queried user exists in db
*/
func TestGetUserSuccess(t *testing.T) {
	mockRepo := new(MockRepository)
	dummyUser := bson.M{
		"_id":         "61098bdd425f4ec901e268df",
		"address":     "Singapore",
		"createdAt":   "1100",
		"description": "I'm young",
		"dob":         "1 Jan 2020",
		"latitude":    9.9,
		"longitude":   11.11,
		"name":        "bean"}
	mockRepo.On("GetUser").Return(dummyUser, nil)
	testService := NewUserService(mockRepo)
	userId := "61098bdd425f4ec901e268df"
	userResult, err := testService.GetUserById(userId)
	assert.Nil(t, err)
	assert.Equal(t, dummyUser, userResult)
}

/**
Testing: GetUserById
Result: Failure as queried user does not exist in db
*/
func TestGetUserFailure(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("GetUser").Return(bson.M{}, errors.New("[Get User] User does not exist in the database"))
	testService := NewUserService(mockRepo)
	userId := "61098bdd425f4ec901e268df"
	userResult, err := testService.GetUserById(userId)
	assert.Nil(t, userResult)
	assert.NotNil(t, err)
	assert.Equal(t, "[Get User] User does not exist in the database", err.Error())
}

// Test deleting an existing user

// Test deleting a user who doesn't exsit -> should throw an error
