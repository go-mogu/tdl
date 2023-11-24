package folders

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/iyear/tdl/app/internal/tgc"
	"golang.org/x/time/rate"
	"time"
)

func Get(ctx context.Context) error {
	c, kvd, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
	fmt.Println(kvd)
	if err != nil {
		return err
	}
	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) (err error) {
		//获取会话文件夹
		filters, err := c.API().MessagesGetDialogFilters(ctx)
		if err != nil {
			return err
		}
		g.Dump(filters)
		return nil
	})

}
