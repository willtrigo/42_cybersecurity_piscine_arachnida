// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   registry.go                                        :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 05:38:49 by dande-je          #+#    #+#             //
//   Updated: 2026/08/26 08:53:00 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type Registry struct {
	readers map[domain.Format]application.MetadataReader
}

func NewRegistry() *Registry {
	jpeg := NewJPEGParser()

	return &Registry{
		readers: map[domain.Format]application.MetadataReader{
			domain.FormatJPEG: jpeg,
		},
	}
}
