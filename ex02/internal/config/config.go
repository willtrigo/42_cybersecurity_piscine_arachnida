// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   config.go                                          :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/25 15:41:05 by dande-je          #+#    #+#             //
//   Updated: 2026/08/25 15:56:53 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package config

import (
	"errors"
	"flag"
	"fmt"
)

var (
	ErrMissingFiles = errors.New("config: at least oner file is required")
)

type Config struct {
	Files []string
}

func Parse(progName string, args []string) (*Config, error) {
	fs := flag.NewFlagSet(progName, flag.ContinueOnError)
	fs.Usage = func() {
		_, _ = fmt.Fprintf(fs.Output(), "Usage: %s FILE1 [FILE2 ...]\n\n", progName)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	files := fs.Args()
	if len(files) == 0 {
		fs.Usage()
		return nil, ErrMissingFiles
	}

	cfg := &Config{Files: files}

	return cfg, nil
}
