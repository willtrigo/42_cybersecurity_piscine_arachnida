// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   metadata.go                                        :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/29 18:36:36 by dande-je          #+#    #+#             //
//   Updated: 2026/09/01 12:47:07 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package format

import (
	"fmt"
	"os"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
)

const timestampLayout = "2006-01-02 15:04:05 MST"

func RenderConsole(result application.InspectionResult) {
	m := result.Metadata

	writeLine("%s\n", m.Path)
	writeLine("  format:\t%s\n", m.Format)
	writeLine("  size:\t\t%d bytes\n", m.Size)
	writeLine("  modified:\t%s\n", m.ModTime.Format("2006-01-02 15:04:05 MST"))
	writeLine("  dimensions:\t%dx%d\n", m.Dimensions.Width, m.Dimensions.Height)

	if !m.HasTags() {
		writeLine("  metadata:\tnone found\n")
	} else {
		writeLine("  metadata:\n")
		for _, tag := range m.Tags {
			var line strings.Builder
			if tag.IDFPath != "" {
				buildLine(&line, "    [%s] %s: %s\n", tag.IDFPath, tag.Name, tag.Value)
			} else {
				buildLine(&line, "    %s: %s\n", tag.Name, tag.Value)
			}
			_, _ = os.Stdout.WriteString(line.String())
		}
	}
}

func RenderGUI(result application.InspectionResult) []string {
	var firstBlock strings.Builder
	var sencodBlock strings.Builder

	m := result.Metadata

	buildLine(&firstBlock, "%s:|:#title\n", m.Path)
	buildLine(&firstBlock, "Format:|:%s\n", m.Format)
	buildLine(&firstBlock, "Size:|:%d bytes\n", m.Size)
	buildLine(&firstBlock, "Modified:|:%s\n", m.ModTime.Format(timestampLayout))
	buildLine(&firstBlock, "Dimensions:|:%d x %d\n", m.Dimensions.Width, m.Dimensions.Height)

	if !m.HasTags() {
		buildLine(&sencodBlock, "Metadata:|:none found\n")
	} else {
		buildLine(&sencodBlock, "Metadata:|:#title\n")
		for _, tag := range m.Tags {
			if tag.IDFPath != "" {
				buildLine(&sencodBlock, "[%s] %s:|:%s\n", tag.IDFPath, tag.Name, tag.Value)
			} else {
				buildLine(&sencodBlock, "%s:|:%s\n", tag.Name, tag.Value)
			}
		}
	}

	return []string{firstBlock.String(), sencodBlock.String()}
}

func writeLine(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}

func buildLine(b *strings.Builder, format string, args ...any) {
	_, _ = fmt.Fprintf(b, format, args...)
}
