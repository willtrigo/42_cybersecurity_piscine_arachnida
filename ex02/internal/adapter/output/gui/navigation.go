// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   navigation.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/30 19:53:21 by dande-je          #+#    #+#             //
//   Updated: 2026/08/30 21:58:09 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
)

const imagePanelRatio = 0.77

type imageViewer struct {
	window fyne.Window

	imageSlot    *fyne.Container
	metadataSlot *fyne.Container
	prevButton   *widget.Button
	nextButton   *widget.Button
	split        *container.Split

	animator *animatedImage

	results []application.InspectionResult
	idx     int
}

func newImageViewer(window fyne.Window, results []application.InspectionResult) (*imageViewer, error) {
	v := &imageViewer{
		window:       window,
		results:      results,
		imageSlot:    container.NewStack(),
		metadataSlot: container.NewStack(),
	}

	v.prevButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), v.showPrevious)
	v.nextButton = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), v.showNext)

	imageArea := container.NewStack(v.imageSlot, newNavigationOverlay(v.prevButton, v.nextButton))

	v.split = container.NewHSplit(imageArea, v.metadataSlot)
	v.split.Offset = imagePanelRatio

	if err := v.show(0); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *imageViewer) Content() fyne.CanvasObject {
	return v.split
}

func (v *imageViewer) showPrevious() {
	if v.idx == 0 {
		return
	}
	if err := v.show(v.idx - 1); err != nil {
		dialog.ShowError(err, v.window)
	}
}

func (v *imageViewer) showNext() {
	if v.idx >= len(v.results)-1 {
		return
	}
	if err := v.show(v.idx + 1); err != nil {
		dialog.ShowError(err, v.window)
	}
}

func (v *imageViewer) show(idx int) error {
	if v.animator != nil {
		v.animator.Stop()
		v.animator = nil
	}

	result := v.results[idx]

	imagePanel, animator, err := newImagePanel(result.Metadata.Path, result.Metadata.Format, result.Metadata.Tags)
	if err != nil {
		return fmt.Errorf("%s: %w", result.Path, err)
	}

	v.idx = idx
	v.animator = animator

	v.imageSlot.Objects = []fyne.CanvasObject{imagePanel}
	v.imageSlot.Refresh()

	v.metadataSlot.Objects = []fyne.CanvasObject{newMetadataPanel(result)}
	v.metadataSlot.Refresh()

	v.updateNavigationState()
	v.window.SetTitle(windowTitle + " - " + result.Metadata.Path)

	if v.animator != nil {
		v.animator.Start()
	}

	return nil
}

func (v *imageViewer) updateNavigationState() {
	if v.idx > 0 {
		v.prevButton.Enable()
	} else {
		v.prevButton.Disable()
	}

	if v.idx < len(v.results)-1 {
		v.nextButton.Enable()
	} else {
		v.nextButton.Disable()
	}
}

func newNavigationOverlay(prev, next *widget.Button) fyne.CanvasObject {
	bg := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 180})
	bg.CornerRadius = 12

	buttons := container.NewHBox(prev, next)

	navBar := container.NewStack(
		bg,
		container.NewPadded(container.NewPadded(buttons)),
	)

	navWrapper := container.NewHBox(
		container.NewPadded(navBar),
		layout.NewSpacer(),
	)

	return container.NewVBox(
		layout.NewSpacer(),
		navWrapper,
	)
}

func (v *imageViewer) Close() {
	if v.animator != nil {
		v.animator.Stop()
	}
}
