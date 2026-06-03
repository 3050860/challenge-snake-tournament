package test

import (
	"context"
	"snake-tournament/internal/repository/mocks"
	"snake-tournament/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGamesDatabase_FindAvailableToEnterGame(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockGamesDatabase)
	game := &models.Game{PlayersAmount: 4}

	mockRepo.On("Find", mock.Anything, mock.Anything, game).Return(nil)

	ctx := context.Background()
	result := mockRepo.FindAvailableToEnterGame(ctx, 4, "test-user-id")

	assert.NotNil(t, result)
	assert.Equal(t, 4, result.PlayersAmount)
	mockRepo.AssertExpectations(t)
}

func TestGamesDatabase_FindAvailableToEnterGames(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockGamesDatabase)
	games := []models.Game{{PlayersAmount: 4}, {PlayersAmount: 6}}

	mockRepo.On("FindAll", mock.Anything, mock.Anything, &games).Return(nil)

	ctx := context.Background()
	result := mockRepo.FindAvailableToEnterGames(ctx, "test-user-id")

	assert.Len(t, result, 2)
	assert.Equal(t, 4, result[0].PlayersAmount)
	assert.Equal(t, 6, result[1].PlayersAmount)
	mockRepo.AssertExpectations(t)
}

func TestGamesDatabase_FindGameByIncludeUser(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockGamesDatabase)
	games := []models.Game{{PlayersAmount: 4}}

	mockRepo.On("FindAll", mock.Anything, mock.Anything, &games).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*[]models.Game)
		*arg = games
	})

	ctx := context.Background()
	result, err := mockRepo.FindGameByIncludeUser(ctx, "test-user-id")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 4, result[0].PlayersAmount)
	mockRepo.AssertExpectations(t)
}
