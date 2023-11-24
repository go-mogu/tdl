package folders

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/tg"
	"github.com/iyear/tdl/app/internal/tgc"
	"golang.org/x/time/rate"
	"time"
)

type DelOptions struct {
	FolderID int
}

func Delete(ctx context.Context, opts DelOptions) error {
	c, kvd, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
	fmt.Println(kvd)
	if err != nil {
		return err
	}
	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) (err error) {
		req := &tg.MessagesUpdateDialogFilterRequest{
			ID: opts.FolderID,
		}
		flag, err := c.API().MessagesUpdateDialogFilter(ctx, req)
		if err != nil {
			return err
		}
		if !flag {
			return gerror.New("Folder ID may be incorrect")
		}
		return nil
	})

}
