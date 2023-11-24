package folders

import (
	"context"
	"fmt"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
	"github.com/iyear/tdl/app/internal/tgc"
	"github.com/iyear/tdl/pkg/storage"
	"github.com/iyear/tdl/pkg/utils"
	"golang.org/x/time/rate"
	"time"
)

type AddOptions struct {
	FolderID        int
	Contacts        bool
	NonContacts     bool
	Groups          bool
	Broadcasts      bool
	Bots            bool
	ExcludeMuted    bool
	ExcludeRead     bool
	ExcludeArchived bool
	ID              int
	Title           string
	Emoticon        string
	PinnedPeers     []string
	IncludePeers    []string
	ExcludePeers    []string
}

func Add(ctx context.Context, opts AddOptions) error {
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
		maxId := 0
		for _, filter := range filters {
			switch v := filter.(type) {
			case *tg.DialogFilter: // dialogFilter#7438f7e8
				if v.ID > maxId {
					maxId = v.ID
				}
			case *tg.DialogFilterChatlist: // dialogFilterChatlist#d64a04a8
				if v.ID > maxId {
					maxId = v.ID
				}
			case *tg.DialogFilterDefault: // dialogFilterDefault#363293ae
				if maxId == 0 {
					maxId = 1
				}
			default:
				panic(v)
			}
		}
		manager := peers.Options{Storage: storage.NewPeers(kvd)}.Build(c.API())
		includePeers := make([]tg.InputPeerClass, 0)
		for _, item := range opts.IncludePeers {
			peer, err := utils.Telegram.GetInputPeer(ctx, manager, item)
			if err != nil {
				continue
			}
			includePeers = append(includePeers, peer.InputPeer())
		}
		req := &tg.MessagesUpdateDialogFilterRequest{
			Flags: 0,
			ID:    maxId + 1,
			Filter: &tg.DialogFilter{
				Title:        opts.Title,
				Emoticon:     opts.Emoticon,
				IncludePeers: includePeers,
			},
		}
		filter, err := c.API().MessagesUpdateDialogFilter(ctx, req)
		if err != nil {
			return err
		}
		fmt.Println(filter)
		return nil
	})

}
