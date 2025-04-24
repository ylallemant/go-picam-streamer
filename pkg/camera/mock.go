//go:build !linux

package camera

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"strings"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/go-picam-streamer/pkg/api"
	"golang.org/x/image/font"
)

const fontPath = "./pkg/camera/UbuntuMono-R.ttf"

func Device(ctx context.Context) (*mock, error) {
	instance := new(mock)

	instance.backgroundColor = color.RGBA{R: 0x30, G: 0x0a, B: 0x24, A: 0xff}
	instance.fontColor = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	instance.fontSize = 32

	ttfBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed load font at: %s", fontPath)
	}

	ttf, err := freetype.ParseFont(ttfBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse font")
	}

	instance.font = ttf

	log.Info().Msgf("stat image generation")
	count := 0
	go func() {
		for {
			log.Info().Msgf("cycle %d", count)
			select {
			case <-ctx.Done():
				return
			default:
				now := time.Now()
				content := now.Format("2006-01-02 15:04:05.999")
				log.Info().Msgf("image %s", content)

				var err error

				fg := image.NewUniform(instance.fontColor)
				bg := image.NewUniform(instance.backgroundColor)

				rgba := image.NewRGBA(image.Rect(0, 0, 1200, 630))
				draw.Draw(rgba, rgba.Bounds(), bg, image.Pt(0, 0), draw.Src)

				text := freetype.NewContext()
				text.SetDPI(72)
				text.SetFont(instance.font)
				text.SetFontSize(instance.fontSize)
				text.SetClip(rgba.Bounds())
				text.SetDst(rgba)
				text.SetSrc(fg)
				text.SetHinting(font.HintingNone)

				textXOffset := 50
				textYOffset := 10 + int(text.PointToFixed(instance.fontSize)>>6) // Note shift/truncate 6 bits first

				pt := freetype.Pt(textXOffset, textYOffset)
				for _, s := range content {
					_, err = text.DrawString(strings.Replace(string(s), "\r", "", -1), pt)
					if err != nil {
						panic(fmt.Sprintf("failed to add text: %s", err.Error()))
					}
					pt.Y += text.PointToFixed(instance.fontSize * 1.5)
				}

				frame := new(bytes.Buffer)
				err = jpeg.Encode(frame, rgba, &jpeg.Options{Quality: 100})
				if err != nil {
					panic(fmt.Sprintf("failed encode jpeg: %s", err.Error()))
				}

				log.Info().Msgf("output image %d", count)
				instance.output <- frame.Bytes()

				time.Sleep(500 * time.Millisecond)
			}
			count = count + 1
		}
	}()

	return instance, nil
}

var _ api.Device = &mock{}

type mock struct {
	output          chan []byte
	font            *truetype.Font
	fontSize        float64
	fontColor       color.RGBA
	backgroundColor color.RGBA
}

func (i *mock) GetOutput() <-chan []byte {
	return i.output
}
