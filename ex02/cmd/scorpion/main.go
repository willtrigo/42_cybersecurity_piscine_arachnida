// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   main.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/21 14:36:07 by dande-je          #+#    #+#             //
//   Updated: 2026/09/02 22:19:15 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package main

import (
	"fmt"
	"os"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/adapter/filesystem"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/adapter/output/cmd"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/adapter/output/gui"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/adapter/parser"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/config"
)

func main() {
	if err := run(os.Args[0], os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "scorpion: %v\n", err)
		os.Exit(1)
	}
}

func run(progName string, args []string) error {
	cfg, err := config.Parse(progName, args)
	if err != nil {
		return err
	}

	registry := parser.NewRegistry()
	inspector := application.NewInspector(registry, filesystem.NewOSSStatReader())

	results, err := inspector.Inspect(cfg.Files)
	if err != nil {
		return err
	}

	if !cfg.GUI {
		presenter := cmd.NewConsolePresenter()
		if err := presenter.Present(results); err != nil {
			return err
		}
	} else {
		presenter := gui.NewWindowPresenter()
		editor := application.NewMetadataEditor(registry)
		if err := presenter.Present(results, editor, inspector); err != nil {
			return err
		}
	}

	return nil
}
