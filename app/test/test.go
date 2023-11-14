package test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/telegram/thumbnail"
	"github.com/gotd/td/tg"
	"github.com/iyear/tdl/app/internal/tgc"
	"golang.org/x/time/rate"
	"image/jpeg"
	"image/png"
	"time"
)

func Test(ctx context.Context) error {
	c, kvd, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
	fmt.Println(kvd)
	if err != nil {
		return err
	}
	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) (err error) {
		dialogFilters, err := c.API().MessagesGetDialogFilters(ctx)
		if err != nil {
			return err
		}
		fmt.Println(dialogFilters)
		self, err := c.Self(ctx)
		if err != nil {
			return err
		}
		photos, err := c.API().PhotosGetUserPhotos(ctx, &tg.PhotosGetUserPhotosRequest{
			UserID: self.AsInput(),
			Limit:  100,
		})
		if err != nil {
			return err
		}
		for _, photoClass := range photos.GetPhotos() {

			photo, b := photoClass.AsNotEmpty()
			if !b {
				return
			}
			var w, h = 100, 100
			realSizes := make([]*tg.PhotoSize, 0)
			var thumb *tg.PhotoStrippedSize
			for _, size := range photo.Sizes {
				switch size.(type) {
				case *tg.PhotoSize:
					photoSize := size.(*tg.PhotoSize)
					realSizes = append(realSizes, photoSize)
				case *tg.PhotoStrippedSize:
					photoSize := size.(*tg.PhotoStrippedSize)
					thumb = photoSize

				}
			}
			if len(realSizes) > 0 {
				w, h = realSizes[len(realSizes)-1].W, realSizes[len(realSizes)-1].H
			}
			fmt.Println(w, h)
			expand, err := thumbnail.Expand(thumb.Bytes)
			if err != nil {
				return err
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
			fmt.Println("IMAGE:" + toString)

		}
		return nil
	})

}
