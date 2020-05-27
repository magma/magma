// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func CheckServiceNameNotExist(ctx context.Context, client *ent.Client, name string) error {
	exist, _ := client.Service.Query().Where(service.Name(name)).Exist(ctx)
	if exist {
		return gqlerror.Errorf("A service with the name %v already exists", name)
	}
	return nil
}

func CheckServiceExternalIDNotExist(ctx context.Context, client *ent.Client, externalID string) error {
	exist, _ := client.Service.Query().Where(service.ExternalID(externalID)).Exist(ctx)
	if exist {
		return gqlerror.Errorf("A service with the external id %v already exists", externalID)
	}
	return nil
}
