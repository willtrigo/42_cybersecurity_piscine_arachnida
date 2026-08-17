// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   filesystem.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/16 21:31:31 by dande-je          #+#    #+#             //
//   Updated: 2026/08/16 21:36:40 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package storage

import (
	"fmt"
	"os"
)

const defaultDirPerm = 0o755

type FilesystemStorage struct {
	basePath string
}

func NewFilesystemStorage(basePath string) (*FilesystemStorage, error) {
	if err := os.MkdirAll(basePath, defaultDirPerm); err != nil {
		return nil, fmt.Errorf("storage: creating output directory %s: %w", basePath, err)
	}

	return &FilesystemStorage{basePath: basePath}, nil
}
