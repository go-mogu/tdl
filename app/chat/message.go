package chat

import (
	"context"
	"fmt"
	pebbledb "github.com/cockroachdb/pebble"
	boltstor "github.com/gotd/contrib/bbolt"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/contrib/pebble"
	"github.com/gotd/contrib/storage"
	"github.com/gotd/td/telegram/query"
	"github.com/gotd/td/telegram/updates"
	updhook "github.com/gotd/td/telegram/updates/hook"
	"github.com/gotd/td/tg"
	"github.com/iyear/tdl/app/internal/tgc"
	"github.com/iyear/tdl/pkg/consts"
	"github.com/iyear/tdl/pkg/logger"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"path/filepath"
	"time"
)

type MsgOptions struct {
	Store bool
}

func Message(ctx context.Context, opts MsgOptions) error {
	if opts.Store {
		return StoreMessage(ctx)
	}
	return NoStoreMessage(ctx)
}

func StoreMessage(ctx context.Context) error {
	log := logger.From(ctx)

	// Peer storage, for resolve caching and short updates handling.
	db, err := pebbledb.Open(filepath.Join(consts.DataDir, "peers.pebble.db"), &pebbledb.Options{})
	if err != nil {
		return errors.Wrap(err, "create pebble storage")
	}
	peerDB := pebble.NewPeerStorage(db)
	log.Info("Storage", zap.String("path", consts.DataDir))

	// Setting up client.
	//
	// Dispatcher is used to register handlers for events.
	dispatcher := tg.NewUpdateDispatcher()
	// Setting up update handler that will fill peer storage before
	// calling dispatcher handlers.
	updateHandler := storage.UpdateHook(dispatcher, peerDB)

	boltdb, err := bbolt.Open(filepath.Join(consts.DataDir, "updates.bolt.db"), 0666, nil)
	if err != nil {
		return errors.Wrap(err, "create bolt storage")
	}
	updatesRecovery := updates.New(updates.Config{
		Handler: updateHandler, // using previous handler with peerDB
		Logger:  log.Named("updates.recovery"),
		Storage: boltstor.NewStateStorage(boltdb),
	})

	// Handler of FLOOD_WAIT that will automatically retry request.
	waiter := floodwait.NewWaiter().WithCallback(func(ctx context.Context, wait floodwait.FloodWait) {
		// Notifying about flood wait.
		log.Warn("Flood wait", zap.Duration("wait", wait.Duration))
		fmt.Println("Got FLOOD_WAIT. Will retry after", wait.Duration)
	})

	// Setup message update handlers.
	dispatcher.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
		log.Info("Channel message", zap.Any("message", update.Message))
		fmt.Println(update.Message)

		return nil
	})

	// Registering handler for new private messages.
	dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateNewMessage) error {
		msg, ok := u.Message.(*tg.Message)
		if !ok {
			return nil
		}
		if msg.Out {
			// Outgoing message.
			return nil
		}

		// Use PeerID to find peer because *Short updates does not contain any entities, so it necessary to
		// store some entities.
		//
		// Storage can be filled using PeerCollector (i.e. fetching all dialogs first).
		p, err := storage.FindPeer(ctx, peerDB, msg.GetPeerID())
		if err != nil {
			return err
		}

		fmt.Printf("%s: %s\n", p, msg.Message)
		return nil
	})

	c, _, err := tgc.NoLogin(ctx, updatesRecovery, ratelimit.New(rate.Every(time.Millisecond*400), 2), waiter)
	if err != nil {
		return err
	}

	return waiter.Run(ctx, func(ctx context.Context) error {
		// Spawning main goroutine.
		if err := tgc.RunWithAuth(ctx, c, func(ctx context.Context) error {
			api := c.API()
			// Getting info about current user.
			self, err := c.Self(ctx)
			if err != nil {
				return errors.Wrap(err, "call self")
			}

			name := self.FirstName
			if self.Username != "" {
				// Username is optional.
				name = fmt.Sprintf("%s (@%s)", name, self.Username)
			}
			fmt.Println("Current user:", name)
			fmt.Println("Filling peer storage from dialogs to cache entities")
			collector := storage.CollectPeers(peerDB)
			if err := collector.Dialogs(ctx, query.GetDialogs(api).Iter()); err != nil {
				return errors.Wrap(err, "collect peers")
			}
			fmt.Println("Filled")

			// Waiting until context is done.
			fmt.Println("Listening for updates. Interrupt (Ctrl+C) to stop.")
			return updatesRecovery.Run(ctx, api, self.ID, updates.AuthOptions{
				IsBot:  self.Bot,
				Forget: false,
				OnStart: func(ctx context.Context) {
					fmt.Println("Update recovery initialized and started, listening for events")
				},
			})
		}); err != nil {
			return errors.Wrap(err, "run")
		}
		return nil
	})
}

func NoStoreMessage(ctx context.Context) error {
	log := logger.From(ctx)

	d := tg.NewUpdateDispatcher()
	gaps := updates.New(updates.Config{
		Handler: d,
		Logger:  log.Named("gaps"),
	})
	d.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		msg, ok := update.Message.(*tg.Message)
		if !ok {
			return nil
		}
		fmt.Printf("msg: %s\n", msg.Message)
		return nil
	})
	c, _, err := tgc.NoLogin(ctx, gaps, ratelimit.New(rate.Every(time.Millisecond*400), 2), updhook.UpdateHook(gaps.Handle))
	if err != nil {
		return err
	}
	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) error {
		// Fetch user info.
		user, err := c.Self(ctx)
		if err != nil {
			return errors.Wrap(err, "call self")
		}

		return gaps.Run(ctx, c.API(), user.ID, updates.AuthOptions{
			OnStart: func(ctx context.Context) {
				log.Info("Gaps started")
			},
		})
	})
}
