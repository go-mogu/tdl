package chat

import (
	"context"
	"fmt"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
	"github.com/iyear/tdl/app/internal/tgc"
	"github.com/iyear/tdl/pkg/logger"
	"github.com/iyear/tdl/pkg/storage"
	"github.com/iyear/tdl/pkg/utils"
	"math/rand"
	"time"

	"github.com/gotd/contrib/middleware/ratelimit"
	msg "github.com/gotd/td/telegram/message"
	"golang.org/x/time/rate"
)

type SendOptions struct {
	Username string
	Msg      string
}

func Send(ctx context.Context, opts SendOptions) error {
	log := logger.From(ctx)

	c, kvd, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
	if err != nil {
		return err
	}

	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) error {
		s := msg.NewSender(c.API())

		manager := peers.Options{Storage: storage.NewPeers(kvd)}.Build(c.API())
		search, err := c.API().ContactsSearch(ctx, &tg.ContactsSearchRequest{
			Q:     opts.Username,
			Limit: 10,
		})

		for _, user := range search.Users {
			user.String()
		}
		peer, err := utils.Telegram.GetInputPeer(ctx, manager, opts.Username)
		if err != nil {
			return fmt.Errorf("failed to get peer: %w", err)
		}

		for i := 0; i < 500; i++ {
			n := rand.Intn(3)
			text, err := s.To(peer.InputPeer()).Text(ctx, fmt.Sprintf("%s: %s", opts.Msg, time.Now().Format("2006-01-02 15:04:05")))
			if err != nil {
				return err
			}
			log.Info(text.String())

			time.Sleep(time.Duration(n+1) * time.Second)
		}

		return nil
	})
}
