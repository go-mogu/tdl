package libtgsconverter

import "bytes"

import "image"
import "image/color"
import "image/gif"

type togif struct {
	gif        gif.GIF
	images     []image.Image
	prev_frame *image.RGBA
}

func (toGif *togif) init(w uint, h uint, options ConverterOptions) {
	toGif.gif.Config.Width = int(w)
	toGif.gif.Config.Height = int(h)
}

func (toGif *togif) SupportsAnimation() bool {
	return true
}

func (toGif *togif) AddFrame(image *image.RGBA, fps uint) error {
	var fpsInt = int(1.0 / float32(fps) * 100.)
	if toGif.prev_frame != nil && sameImage(toGif.prev_frame, image) {
		toGif.gif.Delay[len(toGif.gif.Delay)-1] += fpsInt
		return nil
	}
	toGif.gif.Image = append(toGif.gif.Image, nil)
	toGif.gif.Delay = append(toGif.gif.Delay, fpsInt)
	toGif.gif.Disposal = append(toGif.gif.Disposal, gif.DisposalBackground)
	toGif.images = append(toGif.images, image)
	toGif.prev_frame = image
	return nil
}

func (toGif *togif) Result() []byte {
	q := medianCutQuantizer{mode, nil, false}
	p := q.quantizeMultiple(make([]color.Color, 0, 256), toGif.images)
	// Add transparent entry finally
	var trans_idx uint8 = 0
	if q.reserveTransparent {
		trans_idx = uint8(len(p))
	}
	var id_map = make(map[uint32]uint8)
	for i, img := range toGif.images {
		pi := image.NewPaletted(img.Bounds(), p)
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				c := img.At(x, y)
				cr, cg, cb, ca := c.RGBA()
				cid := (cr>>8)<<16 | cg | (cb >> 8)
				if q.reserveTransparent && ca == 0 {
					pi.Pix[pi.PixOffset(x, y)] = trans_idx
				} else if val, ok := id_map[cid]; ok {
					pi.Pix[pi.PixOffset(x, y)] = val
				} else {
					val := uint8(p.Index(c))
					pi.Pix[pi.PixOffset(x, y)] = val
					id_map[cid] = val
				}
			}
		}
		toGif.gif.Image[i] = pi
	}
	if q.reserveTransparent {
		p = append(p, color.RGBA{0, 0, 0, 0})
	}
	for _, img := range toGif.gif.Image {
		img.Palette = p
	}
	toGif.gif.Config.ColorModel = p
	var data []byte
	w := bytes.NewBuffer(data)
	err := gif.EncodeAll(w, &toGif.gif)
	if err != nil {
		return nil
	}
	return w.Bytes()
}
