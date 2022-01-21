package subscriberdb

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/magma/magma/src/go/protos/magma/subscriberdb"
)

const prefix string = "subdb:"
const prefixPattern string = "subdb:*"

type DB struct {
	redis.Pool
}

type SubscriberDB interface {
	Add(string, subscriberdb.SubscriberData) error
	Delete(string) error
	Update(string, subscriberdb.SubscriberData) error
	Get(string) (*subscriberdb.SubscriberData, error)
	List() ([]*subscriberdb.SubscriberID, error)
}

func NewDB(pool *redis.Pool) *DB {
	return &DB{Pool: *pool}
}

func (db *DB) Add(id string, data subscriberdb.SubscriberData) error {
	serialized, _ := json.Marshal(data)
	if ok, err := db.Pool.Get().Do("SET", prefix+id, string(serialized)); err != nil && ok != true {
		fmt.Println(err)
		return err
	}
	return nil
}

func (db *DB) Delete(id string) error {
	if _, err := db.Pool.Get().Do("DEL", prefix+id); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (db *DB) Update(id string, data subscriberdb.SubscriberData) error {
	db.Add(id, data)
	return nil
}

func (db *DB) Get(id string) (*subscriberdb.SubscriberData, error) {
	data, err := redis.String(db.Pool.Get().Do("GET", prefix+id))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var subData subscriberdb.SubscriberData
	json.Unmarshal([]byte(data), &subData)
	return &subData, nil
}

func (db *DB) List() ([]*subscriberdb.SubscriberID, error) {
	list, err := redis.Strings(db.Pool.Get().Do("KEYS", prefixPattern))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	keys := make([]*subscriberdb.SubscriberID, 0, len(list))
	for _, k := range list {
		keys = append(keys, &subscriberdb.SubscriberID{
			Id: strings.TrimLeft(k, prefix),
		})
	}
	return keys, nil
}
