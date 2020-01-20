// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"

	"github.com/volatiletech/authboss"
	"go.opencensus.io/trace"
)

type traceStorer struct {
	authboss.ServerStorer
	authboss.CreatingServerStorer
	authboss.RememberingServerStorer
}

// TraceStorer wraps a storer and adds tracing to underlying operations.
func TraceStorer(storer authboss.ServerStorer) *traceStorer {
	return &traceStorer{
		ServerStorer:            storer,
		CreatingServerStorer:    authboss.EnsureCanCreate(storer),
		RememberingServerStorer: authboss.EnsureCanRemember(storer),
	}
}

// TraceStatus converts authboss errors to a trace.Status.
func (traceStorer) TraceStatus(err error) trace.Status {
	var code int32
	switch err {
	case nil:
		return trace.Status{Code: trace.StatusCodeOK}
	case authboss.ErrUserFound:
		code = trace.StatusCodeAlreadyExists
	case authboss.ErrUserNotFound, authboss.ErrTokenNotFound:
		code = trace.StatusCodeNotFound
	default:
		code = trace.StatusCodeUnknown
	}
	return trace.Status{Code: code, Message: err.Error()}
}

func (t *traceStorer) Load(ctx context.Context, key string) (authboss.User, error) {
	ctx, span := trace.StartSpan(ctx, "storer.LoadUser")
	span.AddAttributes(trace.StringAttribute("user", key))
	defer span.End()
	user, err := t.ServerStorer.Load(ctx, key)
	span.SetStatus(t.TraceStatus(err))
	return user, err
}

func (t *traceStorer) Save(ctx context.Context, user authboss.User) error {
	ctx, span := trace.StartSpan(ctx, "storer.SaveUser")
	span.AddAttributes(trace.StringAttribute("user", user.GetPID()))
	defer span.End()
	err := t.ServerStorer.Save(ctx, user)
	span.SetStatus(t.TraceStatus(err))
	return err
}

func (t *traceStorer) Create(ctx context.Context, user authboss.User) error {
	ctx, span := trace.StartSpan(ctx, "storer.CreateUser")
	span.AddAttributes(trace.StringAttribute("user", user.GetPID()))
	defer span.End()
	err := t.CreatingServerStorer.Create(ctx, user)
	span.SetStatus(t.TraceStatus(err))
	return err
}

func (t *traceStorer) AddRememberToken(ctx context.Context, pid, token string) error {
	ctx, span := trace.StartSpan(ctx, "storer.AddRememberToken")
	span.AddAttributes(
		trace.StringAttribute("user", pid),
		trace.StringAttribute("token", hashToken(token)),
	)
	defer span.End()
	err := t.RememberingServerStorer.AddRememberToken(ctx, pid, token)
	span.SetStatus(t.TraceStatus(err))
	return err
}

func (t *traceStorer) DelRememberTokens(ctx context.Context, pid string) error {
	ctx, span := trace.StartSpan(ctx, "storer.DelRememberTokens")
	span.AddAttributes(trace.StringAttribute("user", pid))
	defer span.End()
	err := t.RememberingServerStorer.DelRememberTokens(ctx, pid)
	span.SetStatus(t.TraceStatus(err))
	return err
}

func (t *traceStorer) UseRememberToken(ctx context.Context, pid, token string) error {
	ctx, span := trace.StartSpan(ctx, "storer.UseRememberToken")
	span.AddAttributes(
		trace.StringAttribute("user", pid),
		trace.StringAttribute("token", hashToken(token)),
	)
	defer span.End()
	err := t.RememberingServerStorer.UseRememberToken(ctx, pid, token)
	span.SetStatus(t.TraceStatus(err))
	return err
}
