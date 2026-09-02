// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   errors.go                                          :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 16:36:01 by dande-je          #+#    #+#             //
//   Updated: 2026/09/02 20:49:43 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

import "errors"

var (
	ErrUnsupportedFormat = errors.New("domain: unsupported image format")
	ErrUnknownFormat     = errors.New("domain: unrecognized image format")
	ErrWriterUnsupported = errors.New("domain: format does not support writing metadata")
)
