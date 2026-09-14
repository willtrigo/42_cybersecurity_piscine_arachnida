// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   console.go                                         :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/27 08:19:27 by dande-je          #+#    #+#             //
//   Updated: 2026/09/13 12:24:40 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const repeatSize = 80

type ConsolePresenter struct{}

func NewConsolePresenter() *ConsolePresenter {
	return &ConsolePresenter{}
}

func (ConsolePresenter) Present(results []application.InspectionResult) error {
	for i, result := range results {
		renderConsole(result)
		if i < len(results)-1 {
			writeLine("\n%s\n\n", strings.Repeat("-", repeatSize))
		}
	}
	return nil
}

func renderConsole(result application.InspectionResult) {
	m := result.Metadata
	nameWidth := getNameWidth([][]domain.Tag{m.TagsSystem, m.TagsNoneEditable, m.TagsEditable})

	for i, tag := range [][]domain.Tag{m.TagsSystem, m.TagsNoneEditable, m.TagsEditable} {
		if len(tag) > 0 && i > 0 {
			writeLine("%s", "\n")
		}
		for _, t := range tag {
			var line strings.Builder
			buildLine(&line, "%-*s: %s\n", nameWidth, t.Name, t.Value)
			writeLine("%s", line.String())
		}
	}
}

func getNameWidth(tags [][]domain.Tag) int {
	nameWidth := 0
	for _, tag := range tags {
		for _, t := range tag {
			if len(t.Name) > nameWidth {
				nameWidth = len(t.Name) + 1
			}
		}
	}
	return nameWidth
}

func writeLine(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}

func buildLine(b *strings.Builder, format string, args ...any) {
	_, _ = fmt.Fprintf(b, format, args...)
}
