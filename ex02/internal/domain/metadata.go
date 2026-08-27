// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   metadata.go                                        :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 05:45:23 by dande-je          #+#    #+#             //
//   Updated: 2026/08/27 09:58:35 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

import "time"

type Dimensions struct {
	Width  int
	Height int
}

type Tag struct {
	IDFPath string
	Name    string
	Value   string
}

type Metadata struct {
	ModTime    time.Time
	Path       string
	Tags       []Tag
	Format     Format
	Dimensions Dimensions
	Size       int64
}

func (m *Metadata) HasTags() bool {
	return len(m.Tags) > 0
}
