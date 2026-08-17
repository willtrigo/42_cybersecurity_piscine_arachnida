// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   errors.go                                          :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/16 15:08:53 by dande-je          #+#    #+#             //
//   Updated: 2026/08/16 15:21:50 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

import "errors"

var (
	ErrEmptyURL          = errors.New("domain: url must not be empty")
	ErrInvalidURL        = errors.New("domain: url is not a valid absolute url")
	ErrUnsupportedScheme = errors.New("domain: only http and https schemes are supported")
)
