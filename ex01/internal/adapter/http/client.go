// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   client.go                                          :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/17 10:48:30 by dande-je          #+#    #+#             //
//   Updated: 2026/08/17 23:45:47 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package http

import (
	"context"
	"fmt"
	"io"
	nethttp "net/http"
	"time"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/domain"
)

const (
	defaultUserAgent = "spider/1.0 (+42 School Arachnida project; dande-je@student.42sp.org.br)"
	maxResponseBytes = 64 << 20 // 64 MiB
)

const (
	statusMinSuccess = nethttp.StatusOK              // 200
	statusMaxSuccess = nethttp.StatusMultipleChoices // 299
)

type Client struct {
	inner     *nethttp.Client
	userAgent string
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		inner:     &nethttp.Client{Timeout: timeout},
		userAgent: defaultUserAgent,
	}
}

func (c *Client) Get(ctx context.Context, u *domain.URL) ([]byte, string, error) {
	req, err := nethttp.NewRequestWithContext(ctx, nethttp.MethodGet, u.String(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("http: building request for %s: %w", u.String(), err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.inner.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("http: requesting %s: %w", u.String(), err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < statusMinSuccess || resp.StatusCode >= statusMaxSuccess {
		return nil, "", fmt.Errorf("http: %s returned status %d (%s)", u.String(), resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return nil, "", fmt.Errorf("http: reading body of %s: %w", u.String(), err)
	}

	return body, resp.Header.Get("Content-Type"), nil
}
