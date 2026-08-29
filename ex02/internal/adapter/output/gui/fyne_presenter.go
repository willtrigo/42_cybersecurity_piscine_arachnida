// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   fyne_presenter.go                                  :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/28 17:59:02 by dande-je          #+#    #+#             //
//   Updated: 2026/08/28 19:50:34 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
)

const (
	windowTitle  = "Scorpion"
	windowWidth  = 800
	windowHeight = 600
)

type FynePresenter struct {
	app    fyne.App
	window fyne.Window
}

func NewFynePresenter() *FynePresenter {
	fyneApp := app.New()
	window := fyneApp.NewWindow(windowTitle)
	window.Resize(fyne.NewSize(windowWidth, windowHeight))

	return &FynePresenter{app: fyneApp, window: window}
}

func (p *FynePresenter) Present(results []application.InspectionResult) error {
	first, err := firstDisplayableResult(results)
	if err != nil {
		return fmt.Errorf("gui: %w", err)
	}

	image := canvas.NewImageFromFile(first.Metadata.Path)
	image.FillMode = canvas.ImageFillContain

	p.window.SetContent(container.NewStack(image))
	p.window.ShowAndRun()

	return nil
}

func firstDisplayableResult(results []application.InspectionResult) (application.InspectionResult, error) {
	if len(results) == 0 {
		return application.InspectionResult{}, fmt.Errorf("no inspection results to display")
	}

	first := results[0]
	if first.Metadata == nil {
		return application.InspectionResult{}, fmt.Errorf("%s: no metadata available", first.Path)
	}

	return first, nil
}
