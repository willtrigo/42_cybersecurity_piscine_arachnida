// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   config.go                                          :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/16 10:15:06 by dande-je          #+#    #+#             //
//   Updated: 2026/08/17 15:43:34 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

// Package config parses and validates command-line arguments for the
// spider crawler.
package config

import (
	"errors"
	"flag"
	"fmt"
)

const (
	DefaultOutputPath = "./data/"
	DefaultMaxDepth   = 5
)

var (
	// ErrMissingURL is returned when no URL argument is provided.
	ErrMissingURL = errors.New("config: missing required URL argument")
	// ErrInvalidMaxDepth is returned when -l is not a positive integer.
	ErrInvalidMaxDepth = errors.New("config: -l must be a positive integer")
	// ErrMaxDepthWithoutRecur is returned when -l is set without -r.
	ErrMaxDepthWithourRecur = errors.New("config: -l requires -r to be set")
)

// Config holds the parsed and validated command-line configuration.
type Config struct {
	URL        string
	OutputPath string
	MaxDepth   int
	Recursive  bool
}

// Parse parses args into a Config, returning an error if the URL is
// missing or the flag combination is invalid.
func Parse(progName string, args []string) (*Config, error) {
	fs := flag.NewFlagSet(progName, flag.ContinueOnError)
	fs.Usage = func() {
		_, _ = fmt.Fprintf(fs.Output(), "Usage: %s [-r] [-l N] [-p PATH] URL\n\n", progName)
		fs.PrintDefaults()
	}

	recursive := fs.Bool("r", false, "recursively download images found by following links")
	maxDepth := fs.Int("l", DefaultMaxDepth, "maximum recursion depth(only used with -r)")
	outputPath := fs.String("p", DefaultOutputPath, "directory downloaded images are saved to")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return nil, ErrMissingURL
	}

	var maxDepthSet bool
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "l" {
			maxDepthSet = true
		}
	})

	if maxDepthSet && !*recursive {
		return nil, ErrMaxDepthWithourRecur
	}

	if *recursive && *maxDepth <= 0 {
		return nil, fmt.Errorf("%w: got %d", ErrInvalidMaxDepth, *maxDepth)
	}

	cfg := &Config{
		URL:        fs.Arg(0),
		Recursive:  *recursive,
		OutputPath: *outputPath,
	}
	if cfg.Recursive {
		cfg.MaxDepth = *maxDepth
	}

	return cfg, nil
}
