package test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
	"github.com/iyear/tdl/app/internal/tgc"
	"github.com/iyear/tdl/pkg/consts"
	"github.com/iyear/tdl/pkg/storage"
	"github.com/iyear/tdl/pkg/utils"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/time/rate"
	"time"
)

//func Test1(ctx context.Context) error {
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

//func Test1(ctx context.Context) error {
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

//func Test1(ctx context.Context) error {
//	c, kvd, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
//	fmt.Println(kvd)
//	if err != nil {
//		return err
//	}
//	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) (err error) {
//
//		self, err := c.Self(ctx)
//		if err != nil {
//			return err
//		}
//		fmt.Println(self)
//		state, err := c.API().UpdatesGetState(ctx)
//		if err != nil {
//			return err
//		}
//		fmt.Println(state)
//		return nil
//	})
//
//}

func Test1(ctx context.Context) error {
	c, kvd, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
	fmt.Println(kvd)
	if err != nil {
		return err
	}
	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) (err error) {
		photos, err := c.API().PhotosGetUserPhotos(ctx, &tg.PhotosGetUserPhotosRequest{
			UserID: &tg.InputUser{
				UserID:     6404682374,
				AccessHash: 1253764901697455252,
			},
		})
		if err != nil {
			return
		}

		fmt.Println(photos)
		for _, photo := range photos.GetPhotos() {
			d, _ := photo.AsNotEmpty()
			a := tg.InputPhotoFileLocation{
				ID:            d.ID,
				AccessHash:    d.AccessHash,
				FileReference: d.FileReference,
			}

			file, err := c.API().UploadGetFile(ctx, &tg.UploadGetFileRequest{
				Precise:  false,
				Location: &a,
				Limit:    512 * 1024,
			})
			if err != nil {
				return err
			}
			fmt.Println(file)

		}

		return
	})

}

func getLocalCtl(ctx context.Context) *clientv3.Client {
	var localConfig = clientv3.Config{
		Endpoints: []string{"10.8.5.21:2379"},
		Username:  "",
		Password:  "",
	}

	localCtl, err := clientv3.New(localConfig)
	if err != nil {
		g.Log().Fatal(ctx, err)
		return nil
	}
	return localCtl
}

func Test2(ctx context.Context) error {
	ctl := getLocalCtl(ctx)
	nstMap := gmap.NewStrAnyMap()
	getRes, err := ctl.Get(ctx, "/new/tg/", clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range getRes.Kvs {
		key := string(kv.Key)
		g.Log().Info(ctx, "key:", key)
		result, err := gregex.MatchString(`[\d]+`, key)
		if err != nil {
			return err
		}
		ns := result[0]
		g.Log().Info(ctx, "ns:", ns)
		var keys []string
		if nstMap.Contains(ns) {
			list := nstMap.Get(ns)
			keys = list.([]string)
		} else {
			keys = make([]string, 0)
		}
		keys = append(keys, key)
		nstMap.Set(ns, keys)
	}
	var count = 0
	for _, key := range nstMap.Keys() {
		viper.Set(consts.FlagNamespace, key)
		g.Log().Info(ctx, viper.GetString(consts.FlagNamespace))
		c, _, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
		if err != nil {
			return err
		}
		ctx := gctx.New()
		timeout, cancelFunc := context.WithTimeout(ctx, 15*time.Second)
		_ = tgc.RunWithAuth(timeout, c, func(timeout context.Context) (err error) {
			self, err := c.Self(timeout)
			g.Log().Info(timeout, "self", self)
			count++
			return
		})
		cancelFunc()
	}
	g.Log().Info(ctx, "有效账号数量:", count)
	return nil
}

func Test3(ctx context.Context) error {
	c, kvd, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
	fmt.Println(kvd)
	if err != nil {
		return err
	}
	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) (err error) {
		test := gfile.GetContents("/Users/macos/Downloads/test.txt")
		sList := gstr.Split(test, "\r\n")
		strArray := garray.NewStrArrayFrom(sList)
		fmt.Println(strArray)
		count := 0
		manager := peers.Options{Storage: storage.NewPeers(kvd)}.Build(c.API())
		for _, con := range strArray.Range(10, 100) {
			s := message.NewSender(c.API())
			n := grand.N(1, 10)
			msg := `whatsapp filter ,make sales easier, contact:https://t.me/whatsbro1`

			for i := 0; i <= n; i++ {
				msg = "👉" + msg
			}
			peer, err := utils.Telegram.GetInputPeer(ctx, manager, con)
			if err != nil {
				continue
			}
			d := grand.D(3*time.Second, 15*time.Second)
			g.Log().Infof(ctx, "休眠：%v", d)
			time.Sleep(d)
			_, _ = c.API().MessagesReadHistory(ctx, &tg.MessagesReadHistoryRequest{Peer: peer.InputPeer()})
			res, err := s.To(peer.InputPeer()).Text(ctx, msg)
			if err != nil {
				g.Log().Stack(true).Error(ctx, err)
				continue
			}
			g.Log().Info(ctx, res)
			count++
		}
		g.Log().Infof(ctx, "成功发送：%d", count)
		return
	})

}

func Test4(ctx context.Context) error {
	oldPrefix := "/tg/"
	newPrefix := "/new"
	ctl := getLocalCtl(ctx)
	nstMap := gmap.NewStrAnyMap()
	getRes, err := ctl.Get(ctx, oldPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range getRes.Kvs {
		key := string(kv.Key)
		g.Log().Info(ctx, "key:", key)
		result, err := gregex.MatchString(`[\d]+`, key)
		if err != nil {
			return err
		}
		ns := result[0]
		g.Log().Info(ctx, "ns:", ns)
		var keys []string
		if nstMap.Contains(ns) {
			list := nstMap.Get(ns)
			keys = list.([]string)
		} else {
			keys = make([]string, 0)
		}
		keys = append(keys, key)
		nstMap.Set(ns, keys)
	}
	var count = 0
	for _, key := range nstMap.Keys() {
		viper.Set(consts.FlagNamespace, key)
		g.Log().Info(ctx, viper.GetString(consts.FlagNamespace))
		c, _, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
		if err != nil {
			return err
		}
		ctx := gctx.New()
		timeout, cancelFunc := context.WithTimeout(ctx, 30*time.Second)
		err = tgc.RunWithAuth(timeout, c, func(timeout context.Context) (err error) {
			self, err := c.Self(timeout)
			//转移原有session
			g.Log().Info(timeout, "key是否正确:", key == self.Phone, "手机号:", self.Phone)
			g.Log().Info(timeout, "转移原有session", key)
			get, _ := ctl.Get(timeout, oldPrefix+key, clientv3.WithPrefix())
			for _, keyValue := range get.Kvs {
				_, _ = ctl.Put(timeout, newPrefix+gstr.Replace(string(keyValue.Key), key, self.Phone), string(keyValue.Value))
			}
			count++
			return
		})
		if err != nil {
			// session无效，清除session
			g.Log().Warning(ctx, "清理无用session:", key)
			//_, _ = ctl.Delete(ctx, "/tg/"+key, clientv3.WithPrefix())

		}
		cancelFunc()
	}
	g.Log().Info(ctx, "有效账号数量:", count)
	return nil

}

func Test5(ctx context.Context) error {
	c, kvd, err := tgc.NoLogin(ctx, nil, ratelimit.New(rate.Every(time.Millisecond*400), 2))
	fmt.Println(kvd)
	if err != nil {
		return err
	}
	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) (err error) {
		info, err := c.API().UsersGetFullUser(ctx, &tg.InputUserSelf{})
		if err != nil {
			return err
		}
		g.Log().Info(ctx, info)

		return
	})

}
