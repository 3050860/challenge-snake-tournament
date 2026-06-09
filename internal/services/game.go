package services

import (
	"context"
	"fmt"
	"snake-tournament/internal/iface"
	"snake-tournament/models"
	"snake-tournament/models/dto"
	idMutex "snake-tournament/pkg/id_mutex"

	"github.com/sirupsen/logrus"
)

type GameService struct {
	database      iface.IGameRepository
	ticketService iface.ITickets
	editMutex     *idMutex.IdMutex
}

func NewGameService(database iface.IGameRepository, ticketService iface.ITickets) *GameService {
	return &GameService{
		database:      database,
		ticketService: ticketService,
		editMutex:     idMutex.New(),
	}
}

func (s *GameService) Start(ctx context.Context, request dto.GameCreateRequest, user dto.User) (*dto.GameDto, error) {
	logrus.Debugf("Start game by user: %s, game players amount: %d", user.Id, request.PlayersAmount)

	lock := s.editMutex.Lock(user.Id)
	defer lock.Unlock()
	var gameId string
	// var newGame bool
	newGame := false
	game := s.database.FindAvailableToEnterGame(ctx, request.PlayersAmount, user.Id)

	if game == nil {
		logrus.
			WithField("user", user.Email).
			Debug("Games not found, will create new")

		game = models.NewGame(request.PlayersAmount)
		err := s.database.Create(ctx, game)
		if err != nil {
			return nil, fmt.Errorf("failed to create record. error: %w", err)
		}
		gameId = game.GetId()
		newGame = true
	} else {
		gameId = game.GetId()
		logrus.Debugf("Game found: [id=%s, start_time=%s, close_time=%s, new=%s] connect user: %d",
			game.Id, game.StartTime, game.CloseTime, newGame, user.Id)
	}

	// lock = s.editMutex.Lock(gameId)
	err := s.ticketService.CloseTicket(ctx, game, user)
	if err != nil {
		logrus.Debug("Ticket close error")
		if err := s.database.Delete(ctx, game.Id); err != nil {
			//log.Errorf("Delete:%v", game.Id)
		}
		return nil, err
	}
	game.AddPlayer(user)

	err = s.database.Update(ctx, game)
	if err != nil {
		return nil, err
	}

	err = s.database.FindById(ctx, gameId, game)
	if err != nil {
		return nil, err
	}
	// lock.Unlock()

	gameDto := game.ToGameDto()
	return &gameDto, nil
}

func (s *GameService) EnterToGame(ctx context.Context, id string, user dto.User) (*dto.GameDto, error) {
	lock := s.editMutex.Lock(id)

	logrus.Debugf("EnterToGame: connect user: %s to game: %s", user.Id, id)

	var game models.Game
	var gameDto *dto.GameDto
	var err error

	if err := s.database.FindById(ctx, id, &game); err != nil {
		lock.Unlock()
		return nil, err
	}

	if game.IsCloseToEnter(user) {
		logrus.Debugf("EnterToGame %s (user=%s) closed, search new game", game.Id, user.Id)
		gameDto, err = s.Start(ctx, dto.GameCreateRequest{PlayersAmount: game.PlayersAmount}, user)
		if err != nil {
			logrus.Debugf("EnterToGame %s (user=%s) failed to create a new game (%s)", game.Id, user.Id, err.Error())
			lock.Unlock()
			return nil, err
		}
		if game.FullPlayers() {
			// now := time.Now()
			// game.CloseTime = &now
			// s.database.Update(ctx, &game)
			// logrus.Debugf("Game %s set closed", game.Id)
		}
		lock.Unlock()
		logrus.Debugf("EnterToGame %s (user=%s) entered to the game with id %s", game.Id, user.Id, gameDto.Id)
		return gameDto, nil
	}

	if err := s.ticketService.CloseTicket(ctx, &game, user); err != nil {
		//log.Debug("Ticket close error")
		lock.Unlock()
		return nil, err
	}
	// game
	game.AddPlayer(user)
	err = s.database.Update(ctx, &game)
	if err != nil {
		lock.Unlock()
		return nil, err
	}

	//log.Debug("User connected")
	gDto := game.ToGameDto()
	gameDto = &gDto

	lock.Unlock()
	return gameDto, nil
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
