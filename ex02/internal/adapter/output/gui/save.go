// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   save.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/01 20:53:54 by dande-je          #+#    #+#             //
//   Updated: 2026/09/01 21:21:28 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Save struct {
	imageViewer *ImageViewer

	saveButton *widget.Button
	button     fyne.CanvasObject

	idx int
}

func newSave(imageViewer *ImageViewer, idx int) *Save {
	s := &Save{imageViewer: imageViewer, idx: idx}
	s.saveButton = widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), s.saveMetadata)
	s.button = newSaveOverlay(s.saveButton)

	return s
}

func (s *Save) saveMetadata() {
	if err := s.imageViewer.show(s.idx); err != nil {
		dialog.ShowError(err, s.imageViewer.window)
	}
}

func newSaveOverlay(save *widget.Button) fyne.CanvasObject {
	bg := newBg(buttonBgOverlayColorR, buttonBgOverlayColorG, buttonBgOverlayColorB, buttonBgOverlayColorA)
	bg.CornerRadius = cornerRadiusDefault

	navBar := container.NewStack(bg, container.NewPadded(save))

	navWrapper := container.NewHBox(
		layout.NewSpacer(),
		container.NewPadded(navBar),
	)

	return container.NewVBox(layout.NewSpacer(), navWrapper)
}
