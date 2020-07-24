/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package object_store

import (
	"fmt"

	"magma/feg/gateway/registry"
	"magma/orc8r/lib/go/service/config"

	"github.com/go-redis/redis"
)

// RedisClient defines an interface to interact with Redis. Only hash functions
// are used for now
type RedisClient interface {
	HSet(hash string, field string, value string) error
	HGet(hash string, field string) (string, error)
	HGetAll(hash string) (map[string]string, error)
	HDel(hash string, field string) error
}

// RedisClientImpl is the implementation of the redis client using an actual connection
// to redis using go-redis
type RedisClientImpl struct {
	RawClient *redis.Client
}

// NewRedisClient gets the redis configuration from the service config and returns
// a new client or an error if something went wrong
func NewRedisClient() (RedisClient, error) {
	redisConfig, err := config.GetServiceConfig("", registry.REDIS)
	if err != nil {
		return nil, err
	}
	bindAddr, err := redisConfig.GetString("bind")
	if err != nil {
		bindAddr = "127.0.0.1"
	}
	port, err := redisConfig.GetInt("port")
	if err != nil {
		return nil, err
	}
	return &RedisClientImpl{
		RawClient: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", bindAddr, port),
		}),
	}, nil
}

// HSet sets a value at a hash,field pair
func (client *RedisClientImpl) HSet(hash string, field string, value string) error {
	return client.RawClient.HSet(hash, field, value).Err()
}

// HGet gets a value at a hash,field pair
func (client *RedisClientImpl) HGet(hash string, field string) (string, error) {
	return client.RawClient.HGet(hash, field).Result()
}

// HGetAll gets all the possible values for fields in a hash
func (client *RedisClientImpl) HGetAll(hash string) (map[string]string, error) {
	return client.RawClient.HGetAll(hash).Result()
}

// HDel deletes a value at a hash,field pair
func (client *RedisClientImpl) HDel(hash string, field string) error {
	return client.RawClient.HDel(hash, field).Err()
}
