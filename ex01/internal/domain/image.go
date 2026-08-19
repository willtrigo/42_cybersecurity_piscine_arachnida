// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   image.go                                           :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/17 11:04:41 by dande-je          #+#    #+#             //
//   Updated: 2026/08/18 16:23:30 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

import (
	"fmt"
	"path"
	"strings"
)

var SupportExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".gif":  {},
	".bmp":  {},
}

type Image struct {
	URL       *URL
	Filename  string
	Extension string
	Data      []byte
}

func IsSupportedExtension(ext string) bool {
	if ext == "" {
		return false
	}

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	_, ok := SupportExtensions[strings.ToLower(ext)]
	return ok
}

func NewImage(u *URL, data []byte) (*Image, error) {
	if u == nil {
		return nil, ErrEmptyURL
	}
	if len(data) == 0 {
		return nil, ErrEmptyImageData
	}

	filename := path.Base(u.Path())
	extension := strings.ToLower(path.Ext(filename))

	if !IsSupportedExtension(extension) {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedExtension, extension)
	}

	return &Image{
		URL:       u,
		Filename:  filename,
		Extension: extension,
		Data:      data,
	}, nil
}
