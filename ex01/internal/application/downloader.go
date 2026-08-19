// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   downloader.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/17 10:57:15 by dande-je          #+#    #+#             //
//   Updated: 2026/08/18 22:50:13 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/domain"
)

type Downloader struct {
	httpClient domain.HTTPClient
	storage    domain.ImageStorage
}

func NewDownloader(httpClient domain.HTTPClient, storage domain.ImageStorage) *Downloader {
	return &Downloader{httpClient: httpClient, storage: storage}
}

func (d *Downloader) Download(ctx context.Context, u *domain.URL) error {
	if !domain.IsSupportedExtension(extensionOf(u)) {
		return domain.ErrUnsupportedExtension
	}

	body, _, err := d.httpClient.Get(ctx, u)
	if err != nil {
		return err
	}

	img, err := domain.NewImage(u, body)
	if err != nil {
		if errors.Is(err, domain.ErrUnsupportedExtension) {
			return err
		}
		return fmt.Errorf("building image from %s: %w", u.String(), err)
	}

	if err := d.storage.Save(ctx, img); err != nil {
		return err
	}

	return nil
}

func extensionOf(u *domain.URL) string {
	return strings.ToLower(path.Ext(u.Path()))
}
