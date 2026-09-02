// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   editor.go                                          :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/02 18:11:51 by dande-je          #+#    #+#             //
//   Updated: 2026/09/02 18:26:19 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import "github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"

type MetadataEditor struct {
	registry WriterRegistry
}

func NewMetadataEditor(registry WriterRegistry) *MetadataEditor {
	return &MetadataEditor{registry: registry}
}

func (e *MetadataEditor) DeleteTag(path string, format domain.Format, tag string) error {
	writer, err := e.registry.WriterFor(format)
	if err != nil {
		return err
	}
	return writer.DeleteTag(path, tag)
}
