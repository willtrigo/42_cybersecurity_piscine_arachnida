// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   window.go                                          :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/28 17:59:02 by dande-je          #+#    #+#             //
//   Updated: 2026/08/30 18:21:06 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
)

const (
	windowTitle  = "Scorpion"
	windowWidth  = 800
	windowHeight = 500

	imagePanelRatio = 0.77
)

type WindowPresenter struct {
	app    fyne.App
	window fyne.Window
}

func NewWindowPresenter() *WindowPresenter {
	app := app.New()
	window := app.NewWindow(windowTitle)
	window.Resize(fyne.NewSize(windowWidth, windowHeight))

	return &WindowPresenter{app: app, window: window}
}

func (p *WindowPresenter) Present(results []application.InspectionResult) error {
	first := results[0]

	content, animator, err := buildLayout(first)
	if err != nil {
		return err
	}

	if animator != nil {
		p.window.SetOnClosed(animator.Stop)
		animator.Start()
	}

	p.window.SetContent(content)
	p.window.SetTitle("Scorpion - " + first.Metadata.Path)
	p.window.ShowAndRun()

	return nil
}

func buildLayout(result application.InspectionResult) (fyne.CanvasObject, *animatedImage, error) {
	imagePanel, animator, err := newImagePanel(result.Metadata.Path, result.Metadata.Format, result.Metadata.Tags)
	if err != nil {
		return nil, nil, err
	}

	split := container.NewHSplit(imagePanel, newMetadataPanel(result))
	split.Offset = imagePanelRatio

	return split, animator, nil
}
