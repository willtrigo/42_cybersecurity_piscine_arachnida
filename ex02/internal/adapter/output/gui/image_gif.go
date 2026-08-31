// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   image_gif.go                                       :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/30 16:55:08 by dande-je          #+#    #+#             //
//   Updated: 2026/08/31 10:29:06 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	defaultGIFFrameDelay = 100 * time.Millisecond
	gifDelayFactor       = 10
)

type frame struct {
	image image.Image
	delay time.Duration
}

func newAnimatedImagePanel(path string, tags []domain.Tag) (fyne.CanvasObject, *animatedImage, error) {
	frames, err := decodeGIFFrames(path)
	if err != nil {
		return nil, nil, fmt.Errorf("loading animated image: %w", err)
	}

	bg := bgImagePanel()

	if orientation := exifOrientation(tags); orientation != exifOrientationIdentity {
		for i := range frames {
			frames[i].image = applyOrientation(frames[i].image, orientation)
		}
	}

	img := canvas.NewImageFromImage(frames[0].image)
	img.FillMode = canvas.ImageFillContain

	return container.NewStack(bg, container.NewPadded(img)), newAnimatedImage(img, frames), nil
}

func decodeGIFFrames(path string) (frames []frame, err error) {
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("decode: invalid file path")
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, file.Close())
	}()

	decoded, err := gif.DecodeAll(file)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", cleanPath, err)
	}

	screen := image.Rect(0, 0, decoded.Config.Width, decoded.Config.Height)
	canvasImg := image.NewNRGBA(screen)

	result := make([]frame, len(decoded.Image))
	for i, src := range decoded.Image {
		disposal := disposalMethod(decoded, i)

		var preFrameSnapshot *image.NRGBA
		if disposal == gif.DisposalPrevious {
			preFrameSnapshot = cloneNRGBA(canvasImg)
		}

		draw.Draw(canvasImg, src.Bounds(), src, src.Bounds().Min, draw.Over)
		result[i] = frame{image: cloneNRGBA(canvasImg), delay: frameDelay(decoded, i)}

		switch disposal {
		case gif.DisposalBackground:
			draw.Draw(canvasImg, src.Bounds(), image.Transparent, image.Point{}, draw.Src)
		case gif.DisposalPrevious:
			canvasImg = preFrameSnapshot
		}
	}

	return result, nil
}

func disposalMethod(g *gif.GIF, idx int) byte {
	if idx < len(g.Disposal) {
		return g.Disposal[idx]
	}
	return gif.DisposalNone
}

func cloneNRGBA(src *image.NRGBA) *image.NRGBA {
	dst := image.NewNRGBA(src.Bounds())
	copy(dst.Pix, src.Pix)
	return dst
}

func frameDelay(g *gif.GIF, idx int) time.Duration {
	if idx >= len(g.Delay) || g.Delay[idx] <= 0 {
		return defaultGIFFrameDelay
	}
	return time.Duration(g.Delay[idx]) * gifDelayFactor * time.Millisecond
}
