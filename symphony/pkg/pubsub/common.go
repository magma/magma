package pubsub

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/facebookincubator/symphony/pkg/ent"
)

// event log events.
const (
	EntMutation = "ent/mutation"
)

// LogEntry holds an information on a single ent mutation that happened
type LogEntry struct {
	UserName  string    `json:"user_name"`
	UserID    *int      `json:"user_id"`
	Time      time.Time `json:"time"`
	Operation ent.Op    `json:"operation"`
	PrevState *ent.Node `json:"prevState"`
	CurrState *ent.Node `json:"currState"`
}

// Marshal returns the event encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal decodes event data into v.
func Unmarshal(data []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(v)
}
