package tgc

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/gotd/td/telegram"
)

func RunWithAuth(ctx context.Context, client *telegram.Client, f func(ctx context.Context) error) error {
	return client.Run(ctx, func(ctx context.Context) error {
		status, err := client.Auth().Status(ctx)
		if err != nil {
			return err
		}
		if !status.Authorized {
			return fmt.Errorf("not authorized. please login first")
		}
		g.Log().Info(ctx, "Authorized",
			"id", status.User.ID,
			"username", status.User.Username)

		return f(ctx)
	})
}
