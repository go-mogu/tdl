package test

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

//func Test(ctx context.Context) error {
//	c, kvd, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
//	fmt.Println(kvd)
//	if err != nil {
//		return err
//	}
//	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) (err error) {
//		dialogFilters, err := c.API().MessagesGetDialogFilters(ctx)
//		if err != nil {
//			return err
//		}
//		fmt.Println(dialogFilters)
//		self, err := c.Self(ctx)
//		if err != nil {
//			return err
//		}
//		photos, err := c.API().PhotosGetUserPhotos(ctx, &tg.PhotosGetUserPhotosRequest{
//			UserID: self.AsInput(),
//			Limit:  100,
//		})
//		if err != nil {
//			return err
//		}
//		for _, photoClass := range photos.GetPhotos() {
//
//			photo, b := photoClass.AsNotEmpty()
//			if !b {
//				return
//			}
//			var w, h = 100, 100
//			realSizes := make([]*tg.PhotoSize, 0)
//			var thumb *tg.PhotoStrippedSize
//			for _, size := range photo.Sizes {
//				switch size.(type) {
//				case *tg.PhotoSize:
//					photoSize := size.(*tg.PhotoSize)
//					realSizes = append(realSizes, photoSize)
//				case *tg.PhotoStrippedSize:
//					photoSize := size.(*tg.PhotoStrippedSize)
//					thumb = photoSize
//
//				}
//			}
//			if len(realSizes) > 0 {
//				w, h = realSizes[len(realSizes)-1].W, realSizes[len(realSizes)-1].H
//			}
//			fmt.Println(w, h)
//			expand, err := thumbnail.Expand(thumb.Bytes)
//			if err != nil {
//				return err
//			}
//			img, err := jpeg.Decode(bytes.NewReader(expand))
//			if err != nil {
//				panic(err)
//			}
//
//			var buf bytes.Buffer
//			if err := png.Encode(&buf, img); err != nil {
//				panic(err)
//			}
//			toString := base64.StdEncoding.EncodeToString(buf.Bytes())
//			fmt.Println("IMAGE:" + toString)
//
//		}
//		return nil
//	})
//
//}

//func Test(ctx context.Context) error {
//	c, kvd, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
//	fmt.Println(kvd)
//	if err != nil {
//		return err
//	}
//	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) (err error) {
//		reactions, err := c.API().MessagesGetAvailableReactions(ctx, 0)
//		if err != nil {
//			return err
//		}
//
//		fmt.Println(reactions)
//		c.API().PhotosGetUserPhotos(ctx, &tg.PhotosGetUserPhotosRequest{
//			UserID: nil,
//			Offset: 0,
//			MaxID:  0,
//			Limit:  0,
//		})
//		return nil
//	})
//
//}

func Test(ctx context.Context) error {
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
		peer1, _ := utils.Telegram.GetInputPeer(ctx, manager, "it00021hot")
		peer2, err := utils.Telegram.GetInputPeer(ctx, manager, "luxuedng")
		req := &tg.MessagesUpdateDialogFilterRequest{
			Flags: 0,
			ID:    maxId + 1,
			Filter: &tg.DialogFilter{
				Title:        fmt.Sprintf("test%d", maxId),
				IncludePeers: []tg.InputPeerClass{peer1.InputPeer(), peer2.InputPeer()},
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
