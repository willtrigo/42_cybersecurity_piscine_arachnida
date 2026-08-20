// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   main.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/14 19:55:24 by dande-je          #+#    #+#             //
//   Updated: 2026/08/19 23:14:28 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/adapter/http"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/adapter/parser"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/adapter/storage"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/application"
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

	startURL, err := domain.NewURL(cfg.URL)
	if err != nil {
		return fmt.Errorf("invalid URL %q: %w", cfg.URL, err)
	}

	startURL, err = startURL.Normalize()
	if err != nil {
		return fmt.Errorf("normalizing start URL %q: %w", cfg.URL, err)
	}

	fsStorage, err := storage.NewFilesystemStorage(cfg.OutputPath)
	if err != nil {
		return err
	}

	httpClient := http.NewClient(requestTimeout)
	htmlParser := parser.NewHTMLParser()
	downloader := application.NewDownloader(httpClient, fsStorage)
	crawler := application.NewCrawler(httpClient, htmlParser, downloader, cfg.MaxDepth)

	workerCount := runtime.NumCPU() * 2
	crawler.SetWorkerCount(workerCount)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	if err := crawler.Crawl(ctx, startURL); err != nil {
		return fmt.Errorf("\ncrawl failed: %w", err)
	}

	return nil
}
