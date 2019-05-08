/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package cache implements Swx GRPC response cache
package cache

import (
	"container/heap"
	"sync"
	"time"

	"magma/feg/cloud/go/protos"
)

const (
	// defaultTtl - cached entity TTL after last recent use
	DefaultTtl = time.Hour * 3
	// defaultGcInterval - frequency of Garbage Collection checks
	DefaultGcInterval = time.Minute * 5
)

type authEnt struct {
	idx      int
	lastUsed time.Time
	ans      *protos.AuthenticationAnswer
}

type storage struct {
	pq      []*authEnt
	vectors map[string]*authEnt
}

type Impl struct {
	mu   sync.Mutex
	data storage
}

// Go Heap interface implementation
func (s *storage) Len() int {
	return len(s.pq)
}

func (s *storage) Less(i, j int) bool {
	return s.pq[i].lastUsed.Before(s.pq[j].lastUsed)
}

func (s *storage) Swap(i, j int) {
	s.pq[i], s.pq[j] = s.pq[j], s.pq[i]
	s.pq[i].idx = i
	s.pq[j].idx = j
}

func (s *storage) Push(x interface{}) {
	ent := x.(*authEnt)
	ent.idx = len(s.pq)
	s.pq = append(s.pq, ent)
}

func (s *storage) Pop() interface{} {
	n1 := len(s.pq) - 1
	ent := s.pq[n1]
	s.pq = s.pq[0:n1]
	ent.idx = -1
	return ent
}

// New creates & returns a new instance of the cache
func New() *Impl {
	cache, _ := NewExt(DefaultGcInterval, DefaultTtl) // start with default garbage (expired cache) collector
	return cache
}

// NewExt creates & returns a new instance of the cache and GC cancellation chan
func NewExt(interval, ttl time.Duration) (*Impl, chan struct{}) {
	cache := &Impl{data: storage{pq: []*authEnt{}, vectors: map[string]*authEnt{}}}
	return cache, cache.Gc(interval, ttl) // start garbage collector with given interval & ttl
}

// Get retrieves one auth vector from cache if available, adjusts cache and returns the vector, returns nil otherwise
func (swxCache *Impl) Get(imsi string) *protos.AuthenticationAnswer {
	swxCache.mu.Lock()
	defer swxCache.mu.Unlock()
	ent, found := swxCache.data.vectors[imsi]
	if found {
		if len(ent.ans.SipAuthVectors) <= 1 {
			delete(swxCache.data.vectors, imsi)
			heap.Remove(&swxCache.data, ent.idx)
			return ent.ans
		}
		res := *ent.ans // copy answer
		res.SipAuthVectors = res.SipAuthVectors[:1]
		ent.ans.SipAuthVectors = ent.ans.SipAuthVectors[1:]
		ent.lastUsed = time.Now()
		heap.Fix(&swxCache.data, ent.idx)
		return &res
	}
	return nil
}

// Put adds ans vectors into the cache after extracting the first vector from the list, which it returns back to
// the caller in the returned AuthenticationAnswer
func (swxCache *Impl) Put(ans *protos.AuthenticationAnswer) *protos.AuthenticationAnswer {
	if ans == nil || len(ans.UserName) == 0 {
		return ans
	}
	swxCache.mu.Lock()
	defer swxCache.mu.Unlock()
	// delete old cache if present
	ent, found := swxCache.data.vectors[ans.UserName]
	if found {
		heap.Remove(&swxCache.data, ent.idx)
	}
	if len(ans.SipAuthVectors) <= 1 {
		if found {
			delete(swxCache.data.vectors, ans.UserName)
		}
		return ans // only one vector, nothing to cache, just return it
	}

	// cash & return the first vector in a cloned answer
	res := *ans // copy answer
	res.SipAuthVectors = res.SipAuthVectors[:1]
	ans.SipAuthVectors = ans.SipAuthVectors[1:]

	ent = &authEnt{lastUsed: time.Now(), ans: ans}
	swxCache.data.vectors[ans.UserName] = ent
	heap.Push(&swxCache.data, ent)
	return &res
}

// ClearAll removes all cached entities & re-initializes the cache
func (swxCache *Impl) ClearAll() {
	swxCache.mu.Lock()
	swxCache.data = storage{pq: []*authEnt{}, vectors: map[string]*authEnt{}}
	swxCache.mu.Unlock()
}

// Gc starts Garbage Collector with specified check interval & TTL.
// Returns chan to stop the GC: done := cache.Gc(interval, ttl); ... done <- struct{}{}; ...
func (swxCache *Impl) Gc(interval, ttl time.Duration) chan struct{} {
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				stale := t.Add(-ttl)

				swxCache.mu.Lock()
				// Cleanup all expired cache entries
				for swxCache.data.Len() > 0 && swxCache.data.pq[0].lastUsed.Before(stale) {
					delete(swxCache.data.vectors, swxCache.data.pq[0].ans.UserName)
					heap.Pop(&swxCache.data)
				}
				swxCache.mu.Unlock()
			}
		}
	}()
	return done
}
