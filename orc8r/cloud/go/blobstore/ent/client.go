/*
 * Copyright 2020 The Magma Authors
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"log"

	"magma/orc8r/cloud/go/blobstore/ent/migrate"

	"magma/orc8r/cloud/go/blobstore/ent/blob"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Blob is the client for interacting with the Blob builders.
	Blob *BlobClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	c := config{log: log.Println}
	c.options(opts...)
	return &Client{
		config: c,
		Schema: migrate.NewSchema(c.driver),
		Blob:   NewBlobClient(c),
	}
}

// Open opens a connection to the database specified by the driver name and a
// driver-specific data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil

	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %v", err)
	}
	cfg := config{driver: tx, log: c.log, debug: c.debug}
	return &Tx{
		config: cfg,
		Blob:   NewBlobClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Blob.
//		Query().
//		Count(ctx)
//
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := config{driver: dialect.Debug(c.driver, c.log), log: c.log, debug: true}
	return &Client{
		config: cfg,
		Schema: migrate.NewSchema(cfg.driver),
		Blob:   NewBlobClient(cfg),
	}
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// BlobClient is a client for the Blob schema.
type BlobClient struct {
	config
}

// NewBlobClient returns a client for the Blob from the given config.
func NewBlobClient(c config) *BlobClient {
	return &BlobClient{config: c}
}

// Create returns a create builder for Blob.
func (c *BlobClient) Create() *BlobCreate {
	return &BlobCreate{config: c.config}
}

// Update returns an update builder for Blob.
func (c *BlobClient) Update() *BlobUpdate {
	return &BlobUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *BlobClient) UpdateOne(b *Blob) *BlobUpdateOne {
	return c.UpdateOneID(b.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *BlobClient) UpdateOneID(id int) *BlobUpdateOne {
	return &BlobUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Blob.
func (c *BlobClient) Delete() *BlobDelete {
	return &BlobDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *BlobClient) DeleteOne(b *Blob) *BlobDeleteOne {
	return c.DeleteOneID(b.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *BlobClient) DeleteOneID(id int) *BlobDeleteOne {
	return &BlobDeleteOne{c.Delete().Where(blob.ID(id))}
}

// Create returns a query builder for Blob.
func (c *BlobClient) Query() *BlobQuery {
	return &BlobQuery{config: c.config}
}

// Get returns a Blob entity by its id.
func (c *BlobClient) Get(ctx context.Context, id int) (*Blob, error) {
	return c.Query().Where(blob.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *BlobClient) GetX(ctx context.Context, id int) *Blob {
	b, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return b
}
