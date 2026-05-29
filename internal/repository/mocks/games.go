package mocks

import (
	"context"
	"snake-tournament/models"

	"github.com/stretchr/testify/mock"
)

type MockGamesDatabase struct {
	mock.Mock
}

func (m *MockGamesDatabase) FindAvailableToEnterGame(ctx context.Context, playerCount int, userId string) *models.Game {
	args := m.Called(ctx, playerCount, userId)
	return args.Get(0).(*models.Game)
}

func (m *MockGamesDatabase) FindAvailableToEnterGames(ctx context.Context, userId string) []models.Game {
	args := m.Called(ctx, userId)
	return args.Get(0).([]models.Game)
}

func (m *MockGamesDatabase) FindGameByIncludeUser(ctx context.Context, userId string) ([]models.Game, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]models.Game), args.Error(1)
}
