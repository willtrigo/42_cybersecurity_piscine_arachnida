// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   inspector.go                                       :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 16:25:57 by dande-je          #+#    #+#             //
//   Updated: 2026/08/27 08:11:19 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type Inspector struct {
	stat     StatReader
	registry ParserRegistry
}

type InspectionResult struct {
	Metadata *domain.Metadata
	Err      error
	Path     string
}

func NewInspector(registry ParserRegistry, stat StatReader) *Inspector {
	return &Inspector{registry: registry, stat: stat}
}

func (i *Inspector) Inspect(paths []string) ([]InspectionResult, error) {
	results := make([]InspectionResult, len(paths))
	for idx, path := range paths {
		results[idx] = i.inspectOne(path)
	}

	var errs []error
	for _, result := range results {
		if result.Err != nil {
			if errs == nil {
				errs = append(errs, fmt.Errorf("%w", result.Err))
			} else {
				errs = append(errs, fmt.Errorf("scorpion: %w", result.Err))
			}
		}
	}

	if len(errs) > 0 {
		errs = append(errs, fmt.Errorf("scorpion: %d of %d files(s) could not be processed", len(errs), len(results)))
		return results, errors.Join(errs...)
	}

	return results, nil
}

func (i *Inspector) inspectOne(path string) InspectionResult {
	format, err := detectFormat(path)
	if err != nil {
		return InspectionResult{Path: path, Err: fmt.Errorf("%s: %w", path, err)}
	}

	reader, err := i.registry.ReaderFor(format)
	if err != nil {
		return InspectionResult{Path: path, Err: fmt.Errorf("%s: %w", path, err)}
	}

	metadata, err := reader.Read(path)
	if err != nil {
		return InspectionResult{Path: path, Err: fmt.Errorf("%s: %w", path, err)}
	}

	size, modTime, err := i.stat.Stat(path)
	if err != nil {
		return InspectionResult{Path: path, Err: fmt.Errorf("%s: %w", path, err)}
	}

	metadata.Path = path
	metadata.Size = size
	metadata.ModTime = modTime

	return InspectionResult{Path: path, Metadata: metadata}
}

func detectFormat(path string) (domain.Format, error) {
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return domain.FormatUnknown, fmt.Errorf("invalid file path")
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		return domain.FormatUnknown, err
	}
	if info.IsDir() {
		return domain.FormatUnknown, fmt.Errorf("path is a directory")
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return domain.FormatUnknown, err
	}

	if len(data) == 0 {
		return domain.FormatUnknown, fmt.Errorf("empty file")
	}

	headerLen := min(len(data), domain.MaxSignatureLen())

	return domain.DectectFormat(data[:headerLen])
}
