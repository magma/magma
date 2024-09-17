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

package session

import (
	"encoding/json"
	"fbc/cwf/radius/monitoring"
	"fmt"

	"github.com/go-redis/redis"
	"go.opencensus.io/tag"
)

type redisStorage struct {
	data redis.Client
}

func (m *redisStorage) Get(sessionID string) (*State, error) {
	counter := ReadSessionState.Start(
		tag.Upsert(monitoring.SessionIDTag, sessionID),
		tag.Upsert(monitoring.StorageTag, "redis"),
	)

	data, err := m.data.Get(sessionID).Result()

	if err == redis.Nil {
		counter.Failure("not_found")
		return nil, fmt.Errorf("session %s no found in storage", sessionID)
	}

	if err != nil {
		counter.Failure("error")
		return nil, fmt.Errorf("failed to get session %s from storage because of error %v", sessionID, err)
	}

	shapedData := State{}
	err = json.Unmarshal([]byte(data), &shapedData)
	if err != nil {
		counter.Failure("corrupted")
		return nil, ErrInvalidDataFormat
	}

	counter.Success()
	return &shapedData, nil
}

func (m *redisStorage) Set(sessionID string, state State) error {
	counter := WriteSessionState.Start(
		tag.Upsert(monitoring.SessionIDTag, sessionID),
		tag.Upsert(monitoring.StorageTag, "redis"),
	)

	data, err := json.Marshal(state)
	if err != nil {
		counter.Failure("corrupted")
		return ErrInvalidDataFormat
	}

	err = m.data.Set(sessionID, data, 0).Err()
	if err != nil {
		counter.Failure("error")
		return fmt.Errorf("failed to insert session %s from storage because of error %v", sessionID, err)
	}
	counter.Success()
	return nil
}

func (m *redisStorage) Reset(sessionID string) error {
	counter := ResetSessionState.Start(
		tag.Upsert(monitoring.SessionIDTag, sessionID),
		tag.Upsert(monitoring.StorageTag, "redis"),
	)

	err := m.data.Del(sessionID).Err()
	if err != nil {
		counter.Failure("error")
		return fmt.Errorf("failed to reset session %s in redis storage , error is %v", sessionID, err)
	}
	counter.Success()
	return nil
}

// NewMultiSessionRedisStorage Returns a new redis-stored session state storage
func NewMultiSessionRedisStorage(addr string, password string, db int) GlobalStorage {
	return &redisStorage{
		data: *redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}
