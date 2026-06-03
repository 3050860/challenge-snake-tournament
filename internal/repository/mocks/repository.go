package mocks

import (
	"context"
	"snake-tournament/internal/iface"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, model iface.IModel) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *MockRepository) FindAll(ctx context.Context, filter interface{}, models any) error {
	args := m.Called(ctx, filter, models)
	return args.Error(0)
}

func (m *MockRepository) Find(ctx context.Context, filter interface{}, model any) error {
	args := m.Called(ctx, filter, model)
	return args.Error(0)
}

func (m *MockRepository) FindById(ctx context.Context, id string, res any) error {
	args := m.Called(ctx, id, res)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, res iface.IModel) error {
	args := m.Called(ctx, res)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, objectID bson.ObjectID) error {
	args := m.Called(ctx, objectID)
	return args.Error(0)
}

func (m *MockRepository) DeleteAll(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}
