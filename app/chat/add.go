package chat

import (
	"context"
	"fmt"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/tg"
	"github.com/iyear/tdl/app/internal/tgc"
	"golang.org/x/time/rate"
	"time"
)

type AddOptions struct {
	Username string
}

func Add(ctx context.Context, opts AddOptions) error {
	c, _, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
	if err != nil {
		return err
	}

	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) error {
		search, err := c.API().ContactsSearch(ctx, &tg.ContactsSearchRequest{
			Q:     opts.Username,
			Limit: 1,
		})
		if err != nil {
			return err
		}
		if len(search.Users) > 0 {
			for _, user := range search.Users {
				t := user.(*tg.User)
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
		} else {
			importContacts, err := c.API().ContactsImportContacts(ctx, []tg.InputPhoneContact{{Phone: opts.Username}})
			if err != nil {
				return err
			}
			fmt.Println(importContacts.Users)
		}

		return nil
	})
}
