package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"snake-tournament/internal/config"
	"snake-tournament/models"
	"snake-tournament/models/dto"
	"snake-tournament/pkg/ehttp"
	"snake-tournament/pkg/logging"
)

type TicketService struct {
	config *config.Config
}

func NewTicketService(
	config *config.Config,
) *TicketService {
	return &TicketService{
		config: config,
	}
}

func (s TicketService) CloseTicket(ctx context.Context, game models.Game, user dto.User) error {
	logging.GetLogger().Debugf("Close ticket for player amount: %d, user: %s, game id: %s", game.PlayersAmount, user.Id, game.Id)
	ticketCloseDto := dto.CloseTicketRequestDto{
		UserId:       user.Id,
		GameType:     "snake",
		PlayerAmount: game.PlayersAmount,
		GameId:       game.Id,
	}

	requestURL := fmt.Sprintf("http://%s/api/v1/tickets/close-ticket/", s.config.TicketService.Host)
	client := http.Client{}

	jsonBytes, err := json.Marshal(ticketCloseDto)

	if err != nil {
		return fmt.Errorf("failed to marshall data. error: %w", err)
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonBytes))

	if err != nil {
		return err
	}

	req.Header = http.Header{
		"Content-Type": {"application/json"},
	}

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		var errorDto dto.ErrorDto
		err = json.NewDecoder(res.Body).Decode(&errorDto)

		if err != nil {
			return err
		}

		return &ehttp.HttpError{
			Message: errorDto.Message,
			Code:    res.StatusCode,
			Err:     nil,
		}
	}

	return nil
}
