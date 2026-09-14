// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   metadata.go                                        :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 05:45:23 by dande-je          #+#    #+#             //
//   Updated: 2026/09/13 20:38:01 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

type Tag struct {
	IDFPath string
	Name    string
	Value   string
}

type Metadata struct {
	Path             string
	TagsSystem       []Tag
	TagsNoneEditable []Tag
	TagsEditable     []Tag
	Format           Format
}

func (m *Metadata) HasEditableTags() bool {
	return len(m.TagsEditable) > 0
}

func NewTag(idfpath, name, value string) Tag {
	return Tag{IDFPath: idfpath, Name: name, Value: value}
}
