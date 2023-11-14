package chat

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"
	"github.com/iyear/tdl/app/internal/tgc"
	"github.com/iyear/tdl/pkg/libtgsconverter"
	"golang.org/x/image/webp"
	"golang.org/x/time/rate"
	"image/gif"
	"os"
	"sync"
	"time"
)

type StickerOptions struct {
	ShortName string
}

func Download(ctx context.Context, opts StickerOptions) error {
	c, kvd, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
	fmt.Println(kvd)
	if err != nil {
		return err
	}
	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) (err error) {
		set, err := c.API().MessagesGetStickerSet(ctx, &tg.MessagesGetStickerSetRequest{
			Stickerset: &tg.InputStickerSetShortName{
				ShortName: opts.ShortName,
				//ShortName: "UtyaDuck",
				//ShortName: "xiaoshuo",
			},
		})
		fmt.Println(set)
		asModified, _ := set.AsModified()
		wg := sync.WaitGroup{}
		for i, document := range asModified.Documents {
			wg.Add(1)
			document := document
			i := i
			err = grpool.Add(ctx, func(ctx context.Context) {
				defer wg.Done()
				d, _ := document.AsNotEmpty()
				output := new(bytes.Buffer)
				a := tg.InputDocumentFileLocation{
					ID:            d.ID,
					AccessHash:    d.AccessHash,
					FileReference: d.FileReference,
				}
				_, err = downloader.NewDownloader().WithPartSize(512*1024).
					Download(c.API(), &a).
					WithThreads(4).
					Stream(ctx, output)
				if err != nil {
					return
				}
				b := output.Bytes()
				mime := mimetype.Detect(b)
				if mime.Is("image/webp") {
					webpDc, err := webp.Decode(bytes.NewReader(b))
					if err != nil {
						return
					}
					_ = gfile.Mkdir(asModified.Set.ShortName)
					f, err := os.Create(asModified.Set.ShortName + "/" + gconv.String(i) + ".gif")
					if err != nil {
						return
					}
					defer func() {
						if closeErr := f.Close(); err == nil {
							err = closeErr
						}
					}()
					err = gif.Encode(f, webpDc, nil)
					if err != nil {
						return
					}
					//err = gfile.PutBytes(asModified.Set.ShortName+"/"+gconv.String(i)+".webp", b)
					//if err != nil {
					//	return
					//}
				} else {
					opt := libtgsconverter.NewConverterOptions()
					opt.SetExtension("gif")
					//opt.SetScale(0.46875)
					ret, err := libtgsconverter.ImportFromData(b, opt)
					err = gfile.PutBytes(asModified.Set.ShortName+"/"+gconv.String(i)+".gif", ret)
					if err != nil {
						return
					}
					err = gfile.PutBytes(asModified.Set.ShortName+"/"+gconv.String(i)+".tgs", b)
					if err != nil {
						return
					}
				}

			})
			if err != nil {
				return err
			}

		}
		wg.Wait()

		return nil
	})

}
