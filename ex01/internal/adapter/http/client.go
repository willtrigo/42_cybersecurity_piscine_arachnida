// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   client.go                                          :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/17 10:48:30 by dande-je          #+#    #+#             //
//   Updated: 2026/08/17 10:51:08 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package http

import (
	nethttp "net/http"
	"time"
)

type Client struct {
	inner *nethttp.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{inner: &nethttp.Client{Timeout: timeout}}
}
