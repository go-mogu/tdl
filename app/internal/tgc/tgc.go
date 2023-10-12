package tgc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	tdclock "github.com/gotd/td/clock"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/iyear/tdl/pkg/clock"
	"github.com/iyear/tdl/pkg/consts"
	"github.com/iyear/tdl/pkg/key"
	"github.com/iyear/tdl/pkg/kv"
	"github.com/iyear/tdl/pkg/logger"
	"github.com/iyear/tdl/pkg/storage"
	"github.com/iyear/tdl/pkg/utils"
)

func New(ctx context.Context, login bool, updateHandler telegram.UpdateHandler, middlewares ...telegram.Middleware) (*telegram.Client, kv.KV, error) {
	var (
		kvd kv.KV
		err error
	)

	if test := viper.GetString(consts.FlagTest); test != "" {
		kvd, err = kv.NewFile(filepath.Join(os.TempDir(), test)) // persistent storage
	} else {
		kvd, err = kv.NewEtcd(kv.EtcdOptions{
			Ctx:      ctx,
			NS:       viper.GetString(consts.FlagNamespace),
			EndPoint: []string{"10.8.5.21:2379"},
			Username: "",
			Password: "",
		})
	}
	if err != nil {
		return nil, nil, err
	}

	_clock := tdclock.System
	if ntp := viper.GetString(consts.FlagNTP); ntp != "" {
		_clock, err = clock.New()
		if err != nil {
			return nil, nil, err
		}
	}

	mode, err := kvd.Get(key.App())
	if err != nil {
		mode = []byte(consts.AppBuiltin)
	} else {
		if string(mode) != "" {
			login = false
		}
	}
	app, ok := consts.Apps[string(mode)]
	if !ok {
		return nil, nil, fmt.Errorf("can't find app: %s, please try re-login", mode)
	}
	appId, appHash := app.AppID, app.AppHash

	opts := telegram.Options{
		Resolver: dcs.Plain(dcs.PlainOptions{
			Dial: utils.Proxy.GetDial(viper.GetString(consts.FlagProxy)).DialContext,
		}),
		ReconnectionBackoff: func() backoff.BackOff {
			b := backoff.NewExponentialBackOff()

			b.Multiplier = 1.1
			b.MaxElapsedTime = viper.GetDuration(consts.FlagReconnectTimeout)
			b.Clock = _clock
			return b
		},
		Device:         consts.Device,
		SessionStorage: storage.NewSession(kvd, login),
		UpdateHandler:  updateHandler,
		RetryInterval:  5 * time.Second,
		MaxRetries:     5,
		DialTimeout:    10 * time.Second,
		Middlewares:    middlewares,
		Clock:          _clock,
		Logger:         logger.From(ctx).Named("td"),
	}

	// test mode, hook options
	if viper.GetString(consts.FlagTest) != "" {
		appId, appHash = telegram.TestAppID, telegram.TestAppHash
		opts.DC = 2
		opts.DCList = dcs.Test()
		// add rate limit to avoid frequent flood wait
		opts.Middlewares = append(opts.Middlewares, ratelimit.New(rate.Every(100*time.Millisecond), 5))
	}

	logger.From(ctx).Info("New telegram client",
		zap.Int("app", app.AppID),
		zap.String("mode", string(mode)),
		zap.Bool("is_login", login))
	client := telegram.NewClient(appId, appHash, opts)
	return client, kvd, nil
}

func NoLogin(ctx context.Context, updateHandler telegram.UpdateHandler, middlewares ...telegram.Middleware) (*telegram.Client, kv.KV, error) {
	return New(ctx, false, updateHandler, append(middlewares, floodwait.NewSimpleWaiter())...)
}

func Login(ctx context.Context, middlewares ...telegram.Middleware) (*telegram.Client, kv.KV, error) {
	return New(ctx, true, nil, append(middlewares, floodwait.NewSimpleWaiter())...)
}
