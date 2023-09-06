package kv

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Etcd struct {
	ctx context.Context
	ns  string
	kv  *clientv3.Client
}

type EtcdOptions struct {
	Ctx      context.Context
	NS       string   `validate:"required"`
	EndPoint []string `validate:"required"`
	Username string
	Password string
}

func NewEtcd(opts EtcdOptions) (KV, error) {
	cfg := clientv3.Config{
		Context:     opts.Ctx,
		Username:    opts.Username,
		Password:    opts.Password,
		Endpoints:   opts.EndPoint,
		DialTimeout: 15 * time.Second,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Etcd{
		ctx: opts.Ctx,
		ns:  opts.NS,
		kv:  client,
	}, nil
}

func (e *Etcd) Get(key string) ([]byte, error) {
	var val []byte
	resp, err := e.kv.Get(e.ctx, fmt.Sprintf("/%s/%s", e.ns, key))
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, ErrNotFound
	}
	val = resp.Kvs[0].Value
	return val, nil
}

func (e *Etcd) Set(key string, val []byte) error {
	_, err := e.kv.Put(e.ctx, fmt.Sprintf("/%s/%s", e.ns, key), string(val))
	return err
}

// Delete removes a key from the bucket. If the key does not exist then nothing is done and a nil error is returned
func (e *Etcd) Delete(key string) error {
	_, err := e.kv.Delete(e.ctx, fmt.Sprintf("/%s/%s", e.ns, key))
	return err
}
