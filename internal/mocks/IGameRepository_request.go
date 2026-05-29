package mocks

import (
	"snake-tournament/models"

	"github.com/stretchr/testify/mock"
)

func (_mock *MockIGameRepository) FindByIdReturnValues(id string, res models.Game) {
	_mock.On("FindById", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			out := args.Get(2).(*models.Game)
			*out = res
		}).
		Return(nil)
}
