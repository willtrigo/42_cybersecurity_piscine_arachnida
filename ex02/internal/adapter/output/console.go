// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   console.go                                         :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/27 08:19:27 by dande-je          #+#    #+#             //
//   Updated: 2026/08/27 10:39:23 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package output

import (
	"fmt"
	"os"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
)

type ConsolePresenter struct{}

func NewConsolePresenter() *ConsolePresenter {
	return &ConsolePresenter{}
}

func (ConsolePresenter) Present(results []application.InspectionResult) error {
	for _, result := range results {
		if err := presentOne(result); err != nil {
			return err
		}
	}
	return nil
}

func presentOne(result application.InspectionResult) error {
	m := result.Metadata
	if _, err := fmt.Fprintf(os.Stdout, "%s\n", m.Path); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(os.Stdout, "  format:	%s\n", m.Format); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(os.Stdout, "  size:		%d bytes\n", m.Size); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(os.Stdout, "  modified:	%s\n", m.ModTime.Format("2006-01-02 15:04:05 MST")); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(os.Stdout, "  dimensions:	%dx%d\n", m.Dimensions.Width, m.Dimensions.Height); err != nil {
		return err
	}

	if !m.HasTags() {
		_, err := fmt.Fprintf(os.Stdout, "  metadata:	none found\n")
		return err
	}

	if _, err := fmt.Fprintf(os.Stdout, "  metadata:\n"); err != nil {
		return err
	}
	for _, tag := range m.Tags {
		var line string
		if tag.IDFPath != "" {
			line = fmt.Sprintf("    [%s] %s: %s\n", tag.IDFPath, tag.Name, tag.Value)
		} else {
			line = fmt.Sprintf("    %s: %s\n", tag.Name, tag.Value)
		}
		if _, err := fmt.Fprint(os.Stdout, line); err != nil {
			return err
		}
	}

	return nil
}
