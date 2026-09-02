// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   delete.go                                          :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/01 21:30:15 by dande-je          #+#    #+#             //
//   Updated: 2026/09/02 19:53:50 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Delete struct {
	deleteButton *widget.Button
	button       fyne.CanvasObject
	viewer       metadataEditor
	fileName     string
	tag          string
}

func newDelete(fileName, tag string, viewer metadataEditor) *Delete {
	d := &Delete{fileName: fileName, tag: tag, viewer: viewer}
	d.deleteButton = widget.NewButtonWithIcon("", theme.DeleteIcon(), d.deleteMetadata)
	d.button = newDeleteOverlay(d.deleteButton)

	return d
}

func (d *Delete) deleteMetadata() {
	if err := d.viewer.DeleteTag(d.fileName, d.tag); err != nil {
		d.viewer.ShowError(err)
	}
}

func newDeleteOverlay(delete *widget.Button) fyne.CanvasObject {
	bg := newBg(buttonBgOverlayColorR, buttonBgOverlayColorG, buttonBgOverlayColorB, buttonBgOverlayColorA)
	bg.CornerRadius = cornerRadiusDefault

	navBar := container.NewStack(bg, container.NewPadded(delete))

	navWrapper := container.NewHBox(layout.NewSpacer(), container.NewPadded(navBar))

	return container.NewVBox(layout.NewSpacer(), navWrapper)
}
