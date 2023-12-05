package storage

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/gotd/td/telegram"

	"github.com/iyear/tdl/pkg/key"
	"github.com/iyear/tdl/pkg/kv"
)

type Session struct {
	kv    kv.KV
	login bool
}

func NewSession(kv kv.KV, login bool) telegram.SessionStorage {
	return &Session{kv: kv, login: login}
}

func (s *Session) LoadSession(_ context.Context) ([]byte, error) {
	if s.login {
		return nil, nil
	}

	b, err := s.kv.Get(key.Session())
	if err != nil {
		if errors.Is(err, kv.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return b, nil
}

func (s *Session) StoreSession(ctx context.Context, data []byte) error {
	g.Log().Info(ctx, "存储session")
	return s.kv.Set(key.Session(), data)
}
