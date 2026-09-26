// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   errors.go                                          :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 16:36:01 by dande-je          #+#    #+#             //
//   Updated: 2026/09/26 18:29:05 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

import "errors"

var (
	ErrUnsupportedFormat = errors.New("domain: unsupported image format")
	ErrUnknownFormat     = errors.New("domain: unrecognized image format")
	ErrWriterUnsupported = errors.New("domain: format does not support writing metadata")
)

var (
	ErrInvalidSignature   = errors.New("invalid signature")
	ErrTruncatedSignature = errors.New("truncated signature")
)

var (
	ErrTruncatedBMPHeader     = errors.New("truncated DIB header")
	ErrTruncatedBMPFileHeader = errors.New("truncated file header")
)

var (
	ErrMissingIHDRChunk     = errors.New("missing IHDR chunk")
	ErrTruncatedIHDRChunk   = errors.New("truncated IHDR chunk")
	ErrTruncatedChunkHeader = errors.New("truncated chunk header")
	ErrTruncatedChunkType   = errors.New("truncated chunk type")
)
