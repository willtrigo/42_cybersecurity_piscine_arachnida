// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   recover.go                                         :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/24 19:14:22 by dande-je          #+#    #+#             //
//   Updated: 2026/09/24 19:40:25 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import "fmt"

func decodeWithRecover[T any](label string, parse func() T) (result T, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("decode %s header: %w", label, asError(r))
		}
	}()
	result = parse()
	return result, nil
}

func asError(v any) error {
	if err, ok := v.(error); ok {
		return err
	}
	return fmt.Errorf("%v", v)
}
