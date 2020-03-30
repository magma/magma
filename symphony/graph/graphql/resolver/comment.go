// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/ent"
)

type commentResolver struct{}

func (commentResolver) AuthorName(ctx context.Context, obj *ent.Comment) (string, error) {
	author, err := obj.QueryAuthor().Only(ctx)
	if err != nil {
		return "", fmt.Errorf("querying author: %w", err)
	}
	return author.Email, nil
}

func (commentResolver) Author(ctx context.Context, obj *ent.Comment) (*ent.User, error) {
	author, err := obj.QueryAuthor().Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying author: %w", err)
	}
	return author, nil
}
