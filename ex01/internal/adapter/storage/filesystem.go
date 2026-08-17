// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   filesystem.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/16 21:31:31 by dande-je          #+#    #+#             //
//   Updated: 2026/08/17 17:18:58 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/domain"
)

const (
	defaultDirPerm  = 0o755
	defaultFilePerm = 0o644
)

type FilesystemStorage struct {
	basePath string
}

func NewFilesystemStorage(basePath string) (*FilesystemStorage, error) {
	if err := os.MkdirAll(basePath, defaultDirPerm); err != nil {
		return nil, fmt.Errorf("storage: creating output directory %s: %w", basePath, err)
	}

	return &FilesystemStorage{basePath: basePath}, nil
}

func (s *FilesystemStorage) Save(ctx context.Context, img *domain.Image) error {
	destination := filepath.Join(s.basePath, img.Filename)

	if _, err := os.Stat(destination); err == nil {
		return nil
	}

	if err := os.WriteFile(destination, img.Data, defaultFilePerm); err != nil {
		return fmt.Errorf("storage: writing %s: %w", destination, err)
	}

	return nil
}
