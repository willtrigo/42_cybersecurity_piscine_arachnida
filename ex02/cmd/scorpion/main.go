// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   main.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/21 14:36:07 by dande-je          #+#    #+#             //
//   Updated: 2026/08/26 16:54:23 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package main

import (
	"fmt"
	"os"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/adapter/filesystem"
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
	_ = inspector.Inspect(cfg.Files)

	return nil
}
