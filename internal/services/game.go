package services

import (
	"context"
	"fmt"
	"snake-tournament/internal/iface"
	"snake-tournament/models"
	"snake-tournament/models/dto"
	idMutex "snake-tournament/pkg/id_mutex"
	"strconv"

	"github.com/sirupsen/logrus"
)

type GameService struct {
	database          iface.IGameRepository
	ticketService     iface.ITickets
	editMutex         *idMutex.IdMutex
	userMutex         *idMutex.IdMutex
	findOrCreateMutex *idMutex.IdMutex
}

func NewGameService(database iface.IGameRepository, ticketService iface.ITickets) *GameService {
	return &GameService{
		database:          database,
		ticketService:     ticketService,
		editMutex:         idMutex.New(),
		userMutex:         idMutex.New(),
		findOrCreateMutex: idMutex.New(),
	}
}

func (s *GameService) Start(ctx context.Context, request dto.GameCreateRequest, user dto.User) (*dto.GameDto, error) {
	ulock := s.userMutex.Lock(user.Id)
	defer ulock.Unlock()
	return s.startInternal(ctx, request, user)
}

func (s *GameService) startInternal(ctx context.Context, request dto.GameCreateRequest, user dto.User) (*dto.GameDto, error) {
	GameSegment := strconv.Itoa(request.PlayersAmount)
	for attempt := 0; attempt < 10; attempt++ { // ограничение попыток
		// 1. Блокируем поиск/создание
		atom := s.findOrCreateMutex.Lock(GameSegment)

		game := s.database.FindAvailableToEnterGame(ctx, request.PlayersAmount, user.Id)
		var newGame bool

		if game == nil {
			game = models.NewGame(request.PlayersAmount)
			err := s.database.Create(ctx, game)
			if err != nil {
				atom.Unlock()
				return nil, fmt.Errorf("failed to create game: %w", err)
			}
			newGame = true
		}

		gameId := game.GetId()

		// 2. Освобождаем findOrCreateMutex ДО захвата editMutex
		atom.Unlock()

		// 3. Блокируем конкретную игру
		lock := s.editMutex.Lock(gameId)

		// 4. Перечитываем игру ПОСЛЕ захвата блокировки
		if err := s.database.FindById(ctx, gameId, game); err != nil {
			lock.Unlock()
			return nil, err
		}

		// 5. Проверяем, есть ли место
		if len(game.Records) >= game.PlayersAmount {
			lock.Unlock()
			logrus.Debugf("Game %s is full, retrying...", gameId)
			continue // повторить попытку
		}

		// 6. Закрываем тикет
		err := s.ticketService.CloseTicket(ctx, game, user)
		if err != nil {
			lock.Unlock()
			logrus.Debugf("Ticket close failed for game %s, retrying...", gameId)

			// Удаляем новую игру только если в ней никого нет
			if newGame && len(game.Records) == 0 {
				s.database.Delete(ctx, game.Id)
			}
			// continue
			return nil, err
		}

		// 7. Тикет закрыт, добавляем игрока
		game.AddPlayer(user)
		err = s.database.Update(ctx, game)
		if err != nil {
			lock.Unlock()
			return nil, err
		}

		lock.Unlock()
		gameDto := game.ToGameDto()
		return &gameDto, nil
	}

	// Если все попытки не удались (маловероятно, но возможно)
	return nil, fmt.Errorf("failed to join game after multiple attempts")
}

func (s *GameService) EnterToGame(ctx context.Context, id string, user dto.User) (*dto.GameDto, error) {
	ulock := s.userMutex.Lock(user.Id)
	defer ulock.Unlock()

	// Первая проверка без блокировки
	var game models.Game
	if err := s.database.FindById(ctx, id, &game); err != nil {
		return nil, err
	}

	// Если игра закрыта, сразу идём в startInternal
	if game.IsCloseToEnter(user) {
		return s.startInternal(ctx, dto.GameCreateRequest{PlayersAmount: game.PlayersAmount}, user)
	}

	// Игра открыта, пробуем войти
	return s.enterToOpenGameInternal(ctx, id, &game, user)
}

func (s *GameService) enterToOpenGameInternal(ctx context.Context, id string, game *models.Game, user dto.User) (*dto.GameDto, error) {
	for attempt := 0; attempt < 10; attempt++ {
		lock := s.editMutex.Lock(id)

		// Перечитываем игру ПОСЛЕ захвата блокировки
		if err := s.database.FindById(ctx, id, game); err != nil {
			lock.Unlock()
			return nil, err
		}

		// Повторная проверка после блокировки
		if game.IsCloseToEnter(user) {
			lock.Unlock()
			// Игра закрылась, идём в startInternal
			return s.startInternal(ctx, dto.GameCreateRequest{PlayersAmount: game.PlayersAmount}, user)
		}

		// Закрываем тикет
		err := s.ticketService.CloseTicket(ctx, game, user)
		if err != nil {
			lock.Unlock()
			// Тикет не закрылся, идём в startInternal (попробуем другую игру)
			// return s.startInternal(ctx, dto.GameCreateRequest{PlayersAmount: game.PlayersAmount}, user)
			return nil, err
		}

		// Добавляем игрока
		game.AddPlayer(user)
		err = s.database.Update(ctx, game)
		if err != nil {
			lock.Unlock()
			return nil, err
		}

		lock.Unlock()
		gameDto := game.ToGameDto()
		return &gameDto, nil
	}

	return nil, fmt.Errorf("failed to enter game after multiple attempts")
}

func (s *GameService) GetGamesForCurrentUser(ctx context.Context, user dto.User) ([]dto.ResultGameDto, error) {
	//log := logging.GetLogger()
	//log.Debugf("Search games for user: %s", user.Id)

	games, err := s.database.FindGameByIncludeUser(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	response := make([]dto.ResultGameDto, len(games))
	for i := range games {
		response[i] = games[i].ToResultGameDto(user)
	}

	return response, nil
}

func (s *GameService) GetActiveGame(ctx context.Context, user dto.User) ([]dto.GameDto, error) {
	//log := logging.GetLogger()
	//log.Debugf("Search available games for user: %s", user.Id)

	games := s.database.FindAvailableToEnterGames(ctx, user.Id)
	response := make([]dto.GameDto, 0, len(games))
	for i := 0; i < len(games); i++ {
		response = append(response, games[i].ToGameDto())
	}

	return response, nil
}

func (s *GameService) PasteResults(ctx context.Context, id string, request dto.RecordCreateRequest, user dto.User) ([]dto.RecordDto, error) {
	lock := s.editMutex.Lock(id)
	defer lock.Unlock()

	//log := logging.GetLogger()
	//log.Debugf("Paste results to game: %s, by user: %s, result: %d", id, user.Id, request.UserScore)

	var game models.Game
	if err := s.database.FindById(ctx, id, &game); err != nil {
		return nil, err
	}
	logrus.Debugf("FindById:%v", game)

	if err := game.PasteResults(user, request); err != nil {
		return nil, err
	}

	if err := s.database.Update(ctx, &game); err != nil {
		return nil, err
	}

	return game.ToRecordsDto(), nil
}

func (s *GameService) CheckAvailableRenew(ctx context.Context, id string, user dto.User) (bool, error) {
	//log := logging.GetLogger()
	//log.Debugf("Check available renew record to game: %s, by user: %s", id, user.Id)
	var game models.Game
	err := s.database.FindById(ctx, id, &game)
	if err != nil {
		return false, err
	}

	return !game.IsCloseToEnter(user), nil
}

func (s *GameService) SetPrize(ctx context.Context, id string, request dto.SelectPrizeDto, user dto.User) (*dto.ResultGameDto, error) {
	lock := s.editMutex.Lock(id)
	defer lock.Unlock()

	//log := logging.GetLogger()
	//log.Debugf("Set prize to game: %s, by user: %s", id, user.Id)

	var game models.Game
	if err := s.database.FindById(ctx, id, &game); err != nil {
		return nil, err
	}

	if err := game.SetPrizeForUser(user, request.PrizeId, request.Email); err != nil {
		return nil, err
	}

	if err := s.database.Update(ctx, &game); err != nil {
		return nil, err
	}

	response := game.ToResultGameDto(user)
	return &response, nil
}
