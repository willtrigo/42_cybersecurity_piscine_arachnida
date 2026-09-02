// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   registry.go                                        :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 05:38:49 by dande-je          #+#    #+#             //
//   Updated: 2026/09/02 20:29:44 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type Registry struct {
	readers map[domain.Format]application.MetadataReader
	writers map[domain.Format]application.MetadataWriter
}

func NewRegistry() *Registry {
	jpeg := NewJPEGParser()
	png := NewPNGParser()
	gif := NewGIFParser()

	return &Registry{
		readers: map[domain.Format]application.MetadataReader{
			domain.FormatJPEG: jpeg,
			domain.FormatPNG:  png,
			domain.FormatGIF:  gif,
			domain.FormatBMP:  NewBMPParser(),
		},
		writers: map[domain.Format]application.MetadataWriter{
			domain.FormatJPEG: jpeg,
			domain.FormatPNG:  png,
			domain.FormatGIF:  gif,
		},
	}
}

func (r *Registry) ReaderFor(format domain.Format) (application.MetadataReader, error) {
	reader, ok := r.readers[format]
	if !ok {
		return nil, fmt.Errorf("%s, %w", format, domain.ErrUnsupportedFormat)
	}
	return reader, nil
}

func (r *Registry) WriterFor(format domain.Format) (application.MetadataWriter, error) {
	writer, ok := r.writers[format]
	if !ok {
		return nil, fmt.Errorf("%s: %w", format, domain.ErrWriterUnsupported)
	}
	return writer, nil
}
