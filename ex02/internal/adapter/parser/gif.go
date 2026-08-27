// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   gif.go                                             :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/27 11:09:02 by dande-je          #+#    #+#             //
//   Updated: 2026/08/27 15:51:49 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	gifTrailer             = 0x3B
	gifExtensionIntroducer = 0x21
	gifCommentLabel        = 0xFE
	gifImageSeparator      = 0x2C
)

type GIFParser struct{}

func NewGIFParser() *GIFParser {
	return &GIFParser{}
}

func (GIFParser) Read(path string) (*domain.Metadata, error) {
	cleanPath := filepath.Clean(path)

	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("gif: invalid file path")
	}

	data, err := os.Open(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("gif: %w", err)
	}

	r := bufio.NewReader(data)

	header := make([]byte, 6)
	if _, err = io.ReadFull(r, header); err != nil {
		return nil, fmt.Errorf("gif: reading signature: %w", err)
	}

	lsd := make([]byte, 7)
	if _, err = io.ReadFull(r, lsd); err != nil {
		return nil, fmt.Errorf("gif: reading logical screen descriptor: %w", err)
	}
	width := int(binary.LittleEndian.Uint16(lsd[0:2]))
	height := int(binary.LittleEndian.Uint16(lsd[2:4]))
	packed := lsd[4]

	if packed&0x80 != 0 {
		size := globalColorTableSize(packed)
		if _, err = io.CopyN(io.Discard, r, int64(size)); err != nil {
			return nil, fmt.Errorf("gif: skipping global color table: %w", err)
		}
	}

	tags, err := walkBlocks(r)
	if err != nil {
		return nil, fmt.Errorf("gif: %w", err)
	}

	return &domain.Metadata{
		Format:     domain.FormatGIF,
		Dimensions: domain.Dimensions{Width: width, Height: height},
		Tags:       tags,
	}, nil
}

func globalColorTableSize(packed byte) int {
	entries := 1 << ((packed & 0x07) + 1)
	return 3 * entries
}

func walkBlocks(r *bufio.Reader) ([]domain.Tag, error) {
	var tags []domain.Tag

	for {
		introducer, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return tags, nil
			}
			return nil, fmt.Errorf("reading block introducer: %w", err)
		}

		switch introducer {
		case gifTrailer:
			return tags, nil
		case gifExtensionIntroducer:
			label, err := r.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("reading extension label: %w", err)
			}
			data, err := collectSubBlocks(r)
			if err != nil {
				return nil, fmt.Errorf("reading extension body: %w", err)
			}
			if label == gifCommentLabel && len(data) > 0 {
				tags = append(tags, domain.Tag{Name: "Comment", Value: string(data)})
			}
		case gifImageSeparator:
			if err := skipImageBlock(r); err != nil {
				return nil, fmt.Errorf("skipping image block: %w", err)
			}
		default:
			return nil, fmt.Errorf("unexpected block introducer 0x%02X", introducer)
		}
	}
}

func collectSubBlocks(r *bufio.Reader) ([]byte, error) {
	var out []byte
	for {
		size, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if size == 0 {
			return out, nil
		}
		chunk := make([]byte, size)
		if _, err = io.ReadFull(r, chunk); err != nil {
			return nil, err
		}
		out = append(out, chunk...)
	}
}

func skipImageBlock(r *bufio.Reader) error {
	descriptor := make([]byte, 9)
	if _, err := io.ReadFull(r, descriptor); err != nil {
		return fmt.Errorf("reading image descriptor: %w", err)
	}
	packed := descriptor[8]

	if packed&0x80 != 0 {
		size := globalColorTableSize(packed)
		if _, err := io.CopyN(io.Discard, r, int64(size)); err != nil {
			return fmt.Errorf("skipping local color table; %w", err)
		}
	}

	if _, err := r.ReadByte(); err != nil {
		return fmt.Errorf("reading LZW code size: %w", err)
	}

	_, err := collectSubBlocks(r)
	return err
}
