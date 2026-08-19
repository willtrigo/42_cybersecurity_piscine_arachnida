// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   crawler_helpers.go                                 :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/18 09:35:39 by dande-je          #+#    #+#             //
//   Updated: 2026/08/19 00:30:55 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

func (c *Crawler) checkContext(ctx context.Context, msg string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("%s, %w", msg, ctx.Err())
	default:
		return nil
	}
}

func (c *Crawler) isContextError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func isHTMLContent(contentType string) bool {
	return contentType == "" ||
		contentType == "text/html" ||
		contentType == "application/xhtml+xml" ||
		strings.HasPrefix(contentType, "text/html;") ||
		strings.HasPrefix(contentType, "application/xhtml+xml;")
}

func (c *Crawler) logProgess(processed int, state *crawlState, depth int) {
	if processed%10 == 0 {
		c.logger.Printf("spider: processed %d pages, queue size: %d, depth:%d", processed, state.queueSize(), depth)
	}
}
