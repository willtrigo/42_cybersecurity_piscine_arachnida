// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   read.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/03 19:17:47 by dande-je          #+#    #+#             //
//   Updated: 2026/09/10 21:19:25 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type OSFileReader struct{}

func NewOSFileReader() *OSFileReader {
	return &OSFileReader{}
}

func (OSFileReader) ReadFile(path string) ([]byte, error) {
	cleanPath := filepath.Clean(path)

	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("bmp: invalid file path")
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return data, nil
}
