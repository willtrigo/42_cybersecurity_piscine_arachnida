// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   ports.go                                           :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 05:49:30 by dande-je          #+#    #+#             //
//   Updated: 2026/08/26 06:01:09 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import "github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"

type MetadataReader interface {
	Read(path string) (*domain.Metadata, error)
}
