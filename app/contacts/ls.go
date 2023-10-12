package contacts

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/tg"
	"github.com/iyear/tdl/app/internal/tgc"
	"github.com/mattn/go-runewidth"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	"strconv"
	"strings"
	"time"
)

func List(ctx context.Context) error {
	c, _, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
	if err != nil {
		return err
	}

	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) error {
		r, err := c.API().ContactsGetContacts(ctx, 0)

		if err != nil {
			return err
		}
		switch c := r.(type) {
		case *tg.ContactsContacts:
			printTable(c.Users)
		default:
			return errors.Errorf("unexpected type %T", r)
		}
		dc, err := c.API().HelpGetNearestDC(ctx)
		authorization, err := c.API().AuthExportAuthorization(ctx, dc.NearestDC)
		if err != nil {
			return err
		}
		toString := base64.StdEncoding.EncodeToString(authorization.Bytes)
		gfile.PutContents("session.text", toString)
		return nil
	})
}

func printTable(result []tg.UserClass) {
	fmt.Printf("%s %s %s %s %s\n",
		trunc("ID", 10),
		trunc("FirstName", 20),
		trunc("LastName", 20),
		trunc("Username", 20),
		trunc("Phone", 20))

	for _, r := range result {
		user, ok := r.AsNotEmpty()
		if !ok {
			continue
		}
		fmt.Printf("%s %s %s %s %s\n",
			trunc(strconv.FormatInt(r.GetID(), 10), 10),
			trunc(user.FirstName, 20),
			trunc(user.LastName, 20),
			trunc(user.Username, 20),
			trunc(user.Phone, 20))
	}
}

func trunc(s string, len int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		s = "-"
	}
	return runewidth.FillRight(runewidth.Truncate(s, len, "..."), len)
}
