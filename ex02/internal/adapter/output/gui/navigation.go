// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   navigation.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/30 19:53:21 by dande-je          #+#    #+#             //
//   Updated: 2026/09/01 20:38:16 by dande-je         ###   ########.fr       //
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

type Navigation struct {
	imageViewer *ImageViewer

	prevButton *widget.Button
	nextButton *widget.Button

	buttons fyne.CanvasObject

	idx           int
	navigationLen int
}

func newNavigation(imageViewer *ImageViewer, navigationLen int) *Navigation {
	n := &Navigation{imageViewer: imageViewer, idx: 0, navigationLen: navigationLen}
	n.prevButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), n.showPrevious)
	n.nextButton = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), n.showNext)

	n.buttons = newNavigationOverlay(n.prevButton, n.nextButton)
	return n
}

func (n *Navigation) showPrevious() {
	if n.idx == 0 {
		return
	}
	if err := n.imageViewer.show(n.idx - 1); err != nil {
		dialog.ShowError(err, n.imageViewer.window)
	}
}

func (n *Navigation) showNext() {
	if n.idx >= n.navigationLen {
		return
	}
	if err := n.imageViewer.show(n.idx + 1); err != nil {
		dialog.ShowError(err, n.imageViewer.window)
	}
}

func (n *Navigation) updateNavigationState() {
	if n.idx > 0 {
		n.prevButton.Enable()
	} else {
		n.prevButton.Disable()
	}

	if n.idx < n.navigationLen {
		n.nextButton.Enable()
	} else {
		n.nextButton.Disable()
	}
}

func newNavigationOverlay(prev, next *widget.Button) fyne.CanvasObject {
	bg := newBg(buttonBgOverlayColorR, buttonBgOverlayColorG, buttonBgOverlayColorB, buttonBgOverlayColorA)
	bg.CornerRadius = cornerRadiusDefault

	buttons := container.NewHBox(container.NewPadded(prev), container.NewPadded(next))

	navBar := container.NewStack(
		bg,
		container.NewPadded(buttons),
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
