// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   window.go                                          :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/28 17:59:02 by dande-je          #+#    #+#             //
//   Updated: 2026/09/02 19:21:17 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/fstanis/screenresolution"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
)

const (
	windowTitle  = "Scorpion"
	windowWidth  = 800
	windowHeight = 500
)

type WindowPresenter struct {
	app    fyne.App
	window fyne.Window
}

func NewWindowPresenter() *WindowPresenter {
	app := app.New()
	window := app.NewWindow(windowTitle)

	resolution := screenresolution.GetPrimary()
	window.Resize(fyne.NewSize(float32(resolution.Width), float32(resolution.Height)))

	return &WindowPresenter{app: app, window: window}
}

func (p *WindowPresenter) Present(results []application.InspectionResult, editor *application.MetadataEditor) error {
	viewer, err := newImageViewer(p.window, results, editor)
	if err != nil {
		return err
	}

	p.window.SetContent(containerWithMinSize(viewer.Content(), windowWidth, windowHeight))
	p.window.ShowAndRun()

	return nil
}
