// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   crawler_state.go                                   :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/18 09:35:01 by dande-je          #+#    #+#             //
//   Updated: 2026/08/19 23:56:05 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import (
	"context"
	"sync"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/domain"
)

type crawlTask struct {
	url   *domain.URL
	depth int
}

type crawlState struct {
	visited map[string]bool
	start   *domain.URL
	cond    *sync.Cond
	queue   []crawlTask
	mu      sync.Mutex
	pending int
}

func newCrawlState(start *domain.URL) *crawlState {
	s := &crawlState{
		start:   start,
		visited: map[string]bool{start.String(): true},
		queue:   []crawlTask{{url: start, depth: 0}},
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *crawlState) nextTask(ctx context.Context) (task crawlTask, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		if len(s.queue) > 0 {
			task = s.queue[0]
			s.queue = s.queue[1:]
			s.pending++
			return task, true
		}

		if s.pending == 0 || ctx.Err() != nil {
			return crawlTask{}, false
		}

		s.cond.Wait()
	}
}

func (s *crawlState) taskDone() {
	s.mu.Lock()
	s.pending--
	s.mu.Unlock()
	s.cond.Broadcast()
}

func (s *crawlState) addTask(url *domain.URL, depth int) bool {
	s.mu.Lock()
	key := url.String()
	if s.visited[key] {
		s.mu.Unlock()
		return false
	}
	s.visited[url.String()] = true
	s.queue = append(s.queue, crawlTask{url: url, depth: depth})
	s.mu.Unlock()
	s.cond.Broadcast()
	return true
}
