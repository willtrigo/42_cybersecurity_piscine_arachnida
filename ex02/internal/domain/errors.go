// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   errors.go                                          :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 16:36:01 by dande-je          #+#    #+#             //
//   Updated: 2026/08/26 21:05:51 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

import "errors"

var (
	ErrUnsupportedFormat = errors.New("domain: unsupported image format")
	ErrUnknownFormat     = errors.New("doamin: unrecognized image format")
)
