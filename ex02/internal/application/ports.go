// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   ports.go                                           :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 05:49:30 by dande-je          #+#    #+#             //
//   Updated: 2026/09/02 22:38:27 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import (
	"time"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type MetadataReader interface {
	Read(path string) (*domain.Metadata, error)
}

type MetadataWriter interface {
	SetTag() error
	DeleteTag(path string, tag string) error
}

type ParserRegistry interface {
	ReaderFor(format domain.Format) (MetadataReader, error)
}

type WriterRegistry interface {
	WriterFor(format domain.Format) (MetadataWriter, error)
}

type StatReader interface {
	Stat(path string) (size int64, modTime time.Time, err error)
}

type InspectionReload interface {
	InspectionResultReload() ([]InspectionResult, error)
}
