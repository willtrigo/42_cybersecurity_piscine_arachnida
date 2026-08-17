// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   image.go                                           :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/17 11:04:41 by dande-je          #+#    #+#             //
//   Updated: 2026/08/17 15:43:48 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

type Image struct {
	URL       *URL
	Filename  string
	Extension string
	Data      []byte
}
