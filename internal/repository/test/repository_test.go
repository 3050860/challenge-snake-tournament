package test

import (
	"context"
	"snake-tournament/internal/repository/mocks"
	"snake-tournament/models/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MockUser - mock implementation of IModel for testing
type MockUser struct {
	dto.User
}

func (m *MockUser) SetId(id string) {
	m.Id = id
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	model := &MockUser{User: dto.User{Username: "test-user"}}

	mockRepo.On("Create", mock.Anything, model).Return(nil)

	ctx := context.Background()
	err := mockRepo.Create(ctx, model)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindAll(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	var models []dto.User

	mockRepo.On("FindAll", mock.Anything, mock.Anything, &models).Return(nil)

	ctx := context.Background()
	err := mockRepo.FindAll(ctx, bson.M{}, &models)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_Find(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	model := &dto.User{}

	mockRepo.On("Find", mock.Anything, mock.Anything, model).Return(nil)

	ctx := context.Background()
	err := mockRepo.Find(ctx, bson.M{}, model)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindById(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	model := &dto.User{}

	mockRepo.On("FindById", mock.Anything, "test-id", model).Return(nil)

	ctx := context.Background()
	err := mockRepo.FindById(ctx, "test-id", model)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	model := &MockUser{User: dto.User{Id: "test-id"}}

	mockRepo.On("Update", mock.Anything, model).Return(nil)

	ctx := context.Background()
	err := mockRepo.Update(ctx, model)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_Delete(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	objectID, _ := bson.ObjectIDFromHex("507f1f77bcf86cd799439011")

	mockRepo.On("Delete", mock.Anything, objectID).Return(nil)

	ctx := context.Background()
	err := mockRepo.Delete(ctx, objectID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_DeleteAll(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	ids := []string{"507f1f77bcf86cd799439011"}

	mockRepo.On("DeleteAll", mock.Anything, ids).Return(nil)

	ctx := context.Background()
	err := mockRepo.DeleteAll(ctx, ids)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
