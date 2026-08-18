// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   crawler_state.go                                   :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/18 09:35:01 by dande-je          #+#    #+#             //
//   Updated: 2026/08/18 10:28:55 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import "github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/domain"

type crawlTask struct {
	url   *domain.URL
	depth int
}

type crawlState struct {
	visited      map[string]bool
	start        *domain.URL
	currentTask  crawlTask
	queue        []crawlTask
	visitedStart bool
}

func newCrawlState(start *domain.URL) *crawlState {
	return &crawlState{
		start:   start,
		visited: map[string]bool{start.String(): true},
		queue:   []crawlTask{{url: start, depth: 0}},
	}
}

func (s *crawlState) hasMoreTasks() bool {
	return len(s.queue) > 0
}

func (s *crawlState) nextTask() crawlTask {
	task := s.queue[0]
	s.queue = s.queue[1:]
	s.currentTask = task
	return task
}

func (s *crawlState) markPageVisited() {
	s.visitedStart = true
}

func (s *crawlState) hasVisitedStart() bool {
	return s.visitedStart
}

func (s *crawlState) isVisited(url *domain.URL) bool {
	return s.visited[url.String()]
}

func (s *crawlState) addTask(url *domain.URL, depth int) {
	s.visited[url.String()] = true
	s.queue = append(s.queue, crawlTask{url: url, depth: depth})
}

func (s *crawlState) queueSize() int {
	return len(s.queue)
}

func (s *crawlState) currentDepth() int {
	return s.currentTask.depth
}

func (s *crawlState) visitedCount() int {
	return len(s.visited)
}
