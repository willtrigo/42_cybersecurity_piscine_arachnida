// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   image_animation.go                                 :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/30 16:46:10 by dande-je          #+#    #+#             //
//   Updated: 2026/08/30 18:20:55 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type animatedImage struct {
	canvasImage *canvas.Image
	stop        chan struct{}
	frames      []frame
}

func newAnimatedImage(canvasImage *canvas.Image, frames []frame) *animatedImage {
	return &animatedImage{canvasImage: canvasImage, frames: frames, stop: make(chan struct{})}
}

func (a *animatedImage) Start() {
	if len(a.frames) <= 1 {
		return
	}
	go a.loop()
}

func (a *animatedImage) Stop() {
	select {
	case <-a.stop:
	default:
		close(a.stop)
	}
}

func (a *animatedImage) loop() {
	idx := 0
	for {
		current := a.frames[idx]
		fyne.Do(func() {
			a.canvasImage.Image = current.image
			canvas.Refresh(a.canvasImage)
		})

		timer := time.NewTimer(current.delay)
		select {
		case <-a.stop:
			timer.Stop()
			return
		case <-timer.C:
		}

		idx = (idx + 1) % len(a.frames)
	}
}
