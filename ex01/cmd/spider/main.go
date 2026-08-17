// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   main.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/14 19:55:24 by dande-je          #+#    #+#             //
//   Updated: 2026/08/17 15:41:11 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/adapter/http"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/adapter/storage"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/config"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/domain"
)

const requestTimeout = 15 * time.Second

func main() {
	if err := run(os.Args[0], os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "spider: %v\n", err)
		os.Exit(1)
	}
}

func run(progName string, args []string) error {
	cfg, err := config.Parse(progName, args)
	if err != nil {
		return err
	}

	_, err = domain.NewURL(cfg.URL)
	if err != nil {
		return fmt.Errorf("invalid URL %q: %w", cfg.URL, err)
	}

	_, err = storage.NewFilesystemStorage(cfg.OutputPath)
	if err != nil {
		return err
	}

	_ = http.NewClient(requestTimeout)

	return nil
}
