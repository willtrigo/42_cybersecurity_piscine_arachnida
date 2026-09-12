// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   ports.go                                           :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 05:49:30 by dande-je          #+#    #+#             //
//   Updated: 2026/09/12 02:25:07 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import "github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"

type MetadataReader interface {
	Read(path string, stat StatMetadata, file FileReader) (*domain.Metadata, error)
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

type FileReader interface {
	ReadFile(path string) ([]byte, error)
}

type StatReader interface {
	StatIsDir(path string) (bool, error)
	FileStat(path string) (domain.FileStat, error)
}

type StatMetadata interface {
	FileStat(path string) (domain.FileStat, error)
}

type InspectionReload interface {
	InspectionResultReload() ([]InspectionResult, error)
}
