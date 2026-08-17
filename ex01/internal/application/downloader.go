// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   downloader.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/17 10:57:15 by dande-je          #+#    #+#             //
//   Updated: 2026/08/17 17:19:10 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import (
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/domain"
)

type Downloader struct {
	httpClient domain.HTTPClient
	storage    domain.ImageStorage
}

func NewDownloader(httpClient domain.HTTPClient, storage domain.ImageStorage) *Downloader {
	return &Downloader{httpClient: httpClient, storage: storage}
}
