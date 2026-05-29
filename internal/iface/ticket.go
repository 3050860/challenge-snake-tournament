package iface

import (
	"context"

	"github.com/alexsuslov/ehttp/pkg/hticket"
)

type ITickets interface {
	CloseTicket(ctx context.Context, game hticket.IGame, user hticket.IUser) error
}
