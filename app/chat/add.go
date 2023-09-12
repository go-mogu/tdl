package chat

import (
	"context"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/tg"
	"github.com/iyear/tdl/app/internal/tgc"
	"github.com/iyear/tdl/pkg/logger"
	"golang.org/x/time/rate"
	"time"
)

type AddOptions struct {
	Username string
}

func Add(ctx context.Context, opts AddOptions) error {
	log := logger.From(ctx)

	c, _, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
	if err != nil {
		return err
	}

	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) error {
		search, err := c.API().ContactsSearch(ctx, &tg.ContactsSearchRequest{
			Q:     opts.Username,
			Limit: 10,
		})
		if err != nil {
			return err
		}

		//importContacts, err := c.API().ContactsImportContacts(ctx, []tg.InputPhoneContact{{Phone: "+8618028658256"}})
		//if err != nil {
		//	return err
		//}
		//fmt.Println(importContacts.Users)
		for _, user := range search.Users {
			t := user.(*tg.User)
			log.Info(user.String())
			_, err = c.API().ContactsAddContact(ctx, &tg.ContactsAddContactRequest{
				Flags:                    t.Flags,
				AddPhonePrivacyException: true,
				ID:                       t.AsInput(),
				FirstName:                t.FirstName,
				LastName:                 t.LastName,
				Phone:                    t.Phone,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
