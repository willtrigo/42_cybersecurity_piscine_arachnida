// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   inspector.go                                       :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 16:25:57 by dande-je          #+#    #+#             //
//   Updated: 2026/09/12 02:23:55 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type Inspector struct {
	stat     StatReader
	file     FileReader
	registry ParserRegistry
	files    []string
}

type InspectionResult struct {
	Metadata *domain.Metadata
	Err      error
}

func NewInspector(registry ParserRegistry, stat StatReader, file FileReader) *Inspector {
	return &Inspector{registry: registry, stat: stat, file: file}
}

func (i *Inspector) Inspect(paths []string) ([]InspectionResult, error) {
	i.files = paths
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
	format, err := detectFormat(path, i.stat, i.file)
	if err != nil {
		return InspectionResult{Err: fmt.Errorf("%s: %w", path, err)}
	}

	reader, err := i.registry.ReaderFor(format)
	if err != nil {
		return InspectionResult{Err: fmt.Errorf("%s: %w", path, err)}
	}

	metadata, err := reader.Read(path, i.stat, i.file)
	if err != nil {
		return InspectionResult{Err: fmt.Errorf("%s: %w", path, err)}
	}

	return InspectionResult{Metadata: metadata}
}

func detectFormat(path string, stat StatReader, file FileReader) (domain.Format, error) {
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return domain.FormatUnknown, fmt.Errorf("invalid file path")
	}

	isDir, err := stat.StatIsDir(cleanPath)
	if err != nil {
		return domain.FormatUnknown, err
	}
	if isDir {
		return domain.FormatUnknown, fmt.Errorf("path is a directory")
	}

	data, err := file.ReadFile(cleanPath)
	if err != nil {
		return domain.FormatUnknown, err
	}

	if len(data) == 0 {
		return domain.FormatUnknown, fmt.Errorf("empty file")
	}

	headerLen := min(len(data), domain.MaxSignatureLen())

	return domain.DectectFormat(data[:headerLen])
}

func (i *Inspector) InspectionResultReload() ([]InspectionResult, error) {
	newInspectionResult, err := i.Inspect(i.files)
	if err != nil {
		return nil, err
	}
	return newInspectionResult, nil
}
