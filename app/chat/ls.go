package chat

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gotd/td/telegram/thumbnail"
	"image/jpeg"
	"image/png"
	"strconv"
	"strings"
	"time"

	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/telegram/message/peer"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/query"
	"github.com/gotd/td/tg"
	"github.com/mattn/go-runewidth"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/iyear/tdl/app/internal/tgc"
	"github.com/iyear/tdl/pkg/logger"
	"github.com/iyear/tdl/pkg/storage"
	"github.com/iyear/tdl/pkg/texpr"
	"github.com/iyear/tdl/pkg/utils"
)

type Dialog struct {
	ID          int64   `json:"id" comment:"ID of dialog"`
	Type        string  `json:"type" comment:"Type of dialog. Can be 'user', 'channel' or 'group'"`
	Photo       string  `json:"photo"`
	VisibleName string  `json:"visible_name,omitempty" comment:"Title of channel and group, first and last name of user. If empty, output '-'"`
	Username    string  `json:"username,omitempty" comment:"Username of dialog. If empty, output '-'"`
	Topics      []Topic `json:"topics,omitempty" comment:"Topics of dialog. If not set, output '-'"`
}

type Topic struct {
	ID    int    `json:"id" comment:"ID of topic"`
	Title string `json:"title" comment:"Title of topic"`
}

type Output string

var (
	OutputTable Output = "table"
	OutputJSON  Output = "json"
)

// External designation, different from Telegram mtproto
const (
	DialogGroup   = "group"
	DialogPrivate = "private"
	DialogChannel = "channel"
	DialogUnknown = "unknown"
)

type ListOptions struct {
	Output string
	Filter string
}

func List(ctx context.Context, opts ListOptions) error {
	log := logger.From(ctx)

	// align output
	runewidth.EastAsianWidth = false
	runewidth.DefaultCondition.EastAsianWidth = false

	// output available fields
	if opts.Filter == "-" {
		fg := texpr.NewFieldsGetter(nil)
		fields, err := fg.Walk(&Dialog{})
		if err != nil {
			return fmt.Errorf("failed to walk fields: %w", err)
		}

		fmt.Print(fg.Sprint(fields, true))
		return nil
	}
	// compile filter
	filter, err := texpr.Compile(opts.Filter)
	if err != nil {
		return fmt.Errorf("failed to compile filter: %w", err)
	}

	c, kvd, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
	if err != nil {
		return err
	}

	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) error {
		dialogs, err := query.GetDialogs(c.API()).BatchSize(100).Collect(ctx)
		if err != nil {
			return err
		}
		self, err := c.Self(ctx)
		fmt.Println(self)
		blocked, err := utils.Telegram.GetBlockedDialogs(ctx, c.API())
		if err != nil {
			return err
		}
		manager := peers.Options{Storage: storage.NewPeers(kvd)}.Build(c.API())

		result := make([]*Dialog, 0, len(dialogs))
		for _, d := range dialogs {
			id := utils.Telegram.GetInputPeerID(d.Peer)
			// we can update our access hash state if there is any new peer.
			if err = applyPeers(ctx, manager, d.Entities, id); err != nil {
				log.Warn("failed to apply peer updates", zap.Int64("id", id), zap.Error(err))
			}

			// filter blocked peers
			if _, ok := blocked[id]; ok {
				continue
			}
			var r *Dialog
			switch t := d.Peer.(type) {
			case *tg.InputPeerUser:
				r = processUser(t.UserID, d.Entities)
			case *tg.InputPeerChannel:
				r = processChannel(ctx, c.API(), t.ChannelID, d.Entities)
			case *tg.InputPeerChat:
				r = processChat(t.ChatID, d.Entities)
			}

			// skip unsupported types
			if r == nil {
				continue
			}

			// filter
			b, err := texpr.Run(filter, r)
			if err != nil {
				return fmt.Errorf("failed to run filter: %w", err)
			}
			if !b.(bool) {
				continue
			}

			result = append(result, r)
		}

		switch Output(opts.Output) {
		case OutputTable:
			printTable(result)
		case OutputJSON:
			bytes, err := json.MarshalIndent(result, "", "\t")
			if err != nil {
				return fmt.Errorf("marshal json: %w", err)
			}

			fmt.Println(string(bytes))
		default:
			return fmt.Errorf("unknown output: %s", opts.Output)
		}

		return nil
	})
}

func printTable(result []*Dialog) {
	fmt.Printf("%s %s %s %s %s\n",
		trunc("ID", 10),
		trunc("Type", 8),
		trunc("VisibleName", 20),
		trunc("Username", 20),
		"Topics")

	for _, r := range result {
		fmt.Printf("%s %s %s %s %s\n",
			trunc(strconv.FormatInt(r.ID, 10), 10),
			trunc(r.Type, 8),
			trunc(r.VisibleName, 20),
			trunc(r.Username, 20),
			topicsString(r.Topics))
	}
}

func trunc(s string, len int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		s = "-"
	}

	return runewidth.FillRight(runewidth.Truncate(s, len, "..."), len)
}

func topicsString(topics []Topic) string {
	if len(topics) == 0 {
		return "-"
	}

	s := make([]string, 0, len(topics))
	for _, t := range topics {
		s = append(s, fmt.Sprintf("%d: %s", t.ID, t.Title))
	}

	return strings.Join(s, ", ")
}

func processUser(id int64, entities peer.Entities) *Dialog {
	u, ok := entities.User(id)
	if !ok {
		return nil
	}
	photo, b := u.Photo.AsNotEmpty()
	if b {
		expand, err := thumbnail.Expand(photo.StrippedThumb)
		if err != nil {
			panic(err)
		}
		img, err := jpeg.Decode(bytes.NewReader(expand))
		if err != nil {
			panic(err)
		}

		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			panic(err)
		}
		toString := base64.StdEncoding.EncodeToString(buf.Bytes())
		fmt.Println(toString)
	}

	return &Dialog{
		ID:          u.ID,
		VisibleName: visibleName(u.FirstName, u.LastName),
		Username:    u.Username,
		Type:        DialogPrivate,
		Topics:      nil,
	}
}

func processChannel(ctx context.Context, api *tg.Client, id int64, entities peer.Entities) *Dialog {

	c, ok := entities.Channel(id)
	if !ok {
		return nil
	}
	d := &Dialog{
		ID:          c.ID,
		VisibleName: c.Title,
		Username:    c.Username,
	}

	// channel type
	switch {
	case c.Broadcast:
		d.Type = DialogChannel
	case c.Megagroup, c.Gigagroup:
		d.Type = DialogGroup
	default:
		d.Type = DialogUnknown
	}

	if c.Forum {
		req := &tg.ChannelsGetForumTopicsRequest{
			Channel: c.AsInput(),
			Limit:   100,
		}

		topics, err := api.ChannelsGetForumTopics(ctx, req)
		if err != nil {
			return nil
		}

		d.Topics = make([]Topic, 0, len(topics.Topics))
		for _, tp := range topics.Topics {
			if t, ok := tp.(*tg.ForumTopic); ok {
				d.Topics = append(d.Topics, Topic{
					ID:    t.ID,
					Title: t.Title,
				})
			}
		}
	}

	return d
}

func processChat(id int64, entities peer.Entities) *Dialog {
	c, ok := entities.Chat(id)
	if !ok {
		return nil
	}
	return &Dialog{
		ID:          c.ID,
		VisibleName: c.Title,
		Username:    "-",
		Type:        DialogGroup,
		Topics:      nil,
	}
}

func visibleName(first, last string) string {
	if first == "" && last == "" {
		return ""
	}

	if first == "" {
		return last
	}

	if last == "" {
		return first
	}

	return first + " " + last
}

func applyPeers(ctx context.Context, manager *peers.Manager, entities peer.Entities, id int64) error {
	users := make([]tg.UserClass, 0, 1)
	if user, ok := entities.User(id); ok {
		users = append(users, user)
	}

	chats := make([]tg.ChatClass, 0, 1)
	if chat, ok := entities.Chat(id); ok {
		chats = append(chats, chat)
	}

	return manager.Apply(ctx, users, chats)
}
