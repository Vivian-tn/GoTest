// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errgroup provides synchronization, error propagation, and Context
// cancelation for groups of goroutines working on subtasks of a common task.
//
// Copied from golang.org/x/sync/errgroup
package errgroup

import (
	"fmt"
	"runtime/debug"
	"sync"

	"git.in.zhihu.com/zhsearch/search-ingress/pkg/log"
	"git.in.zhihu.com/zhsearch/search-ingress/pkg/util/statsd"
	"golang.org/x/net/context"
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid and does not cancel on error.
type Group struct {
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

const groupMark string = "GROUP.GOTRACE"

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(groupMark, r, string(debug.Stack()))
				if e, ok := r.(error); ok {
					log.WithError(fmt.Errorf("[%s]ingress-panic: %w", groupMark, e))
				} else {
					log.WithError(fmt.Errorf("[%s]ingress-panic: %+v", groupMark, r))
				}
				//debug.PrintStack()

				statsd.Increment("search-ingress.bug.panic.count")
			}
		}()

		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}
