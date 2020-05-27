// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"gocloud.dev/pubsub"
)

// Subscriber represents types than can open subscriptions.
type Subscriber interface {
	Subscribe(context.Context) (*pubsub.Subscription, error)
}

// The SubscriberFunc type is an adapter to allow the use of
// ordinary functions as subscribers.
type SubscriberFunc func(context.Context) (*pubsub.Subscription, error)

// Subscribe returns f(ctx).
func (f SubscriberFunc) Subscribe(ctx context.Context) (*pubsub.Subscription, error) {
	return f(ctx)
}

// URLSubscriber opens subscriptions from urls.
type URLSubscriber string

// NewURLSubscriber creates a url subscriber from a url string.
func NewURLSubscriber(url string) URLSubscriber {
	return URLSubscriber(url)
}

// Subscribe opens a subscription from url.
func (u URLSubscriber) Subscribe(ctx context.Context) (*pubsub.Subscription, error) {
	return pubsub.OpenSubscription(ctx, u.String())
}

// String returns the textual representation of url opener.
func (u URLSubscriber) String() string {
	return string(u)
}

// Set updates the value of the url opener.
func (u *URLSubscriber) Set(v string) error {
	if _, err := url.Parse(v); err != nil {
		return fmt.Errorf("parsing url: %w", err)
	}
	*u = URLSubscriber(v)
	return nil
}

// NewNopSubscriber returns a subscriber that always fails to open subscriptions.
func NewNopSubscriber() Subscriber {
	return SubscriberFunc(func(context.Context) (*pubsub.Subscription, error) {
		return nil, errors.New("nop subscriber")
	})
}
