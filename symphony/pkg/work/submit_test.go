// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package work

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gocloud.dev/pubsub/mempubsub"
)

func TestNewTopicSubmitter(t *testing.T) {
	t.Run("WithTopic", func(t *testing.T) {
		submitter := NewSubmitter(mempubsub.NewTopic())
		assert.NotNil(t, submitter)
	})
	t.Run("WithURL", func(t *testing.T) {
		submitter, err := NewSubmitterURL(
			context.Background(),
			"mem://test-topic",
		)
		assert.NoError(t, err)
		assert.NotNil(t, submitter)
	})
	t.Run("WithBadURL", func(t *testing.T) {
		_, err := NewSubmitterURL(
			context.Background(),
			"bad://test-topic",
		)
		assert.Error(t, err)
	})
}

func TestTopicSubmit(t *testing.T) {
	submitter := NewSubmitter(mempubsub.NewTopic())
	defer func() {
		err := submitter.Close()
		assert.NoError(t, err)
	}()

	subscription := mempubsub.NewSubscription(submitter.topic, time.Second)
	defer func() {
		err := subscription.Shutdown(context.Background())
		assert.NoError(t, err)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	const nr = 128
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < nr; i++ {
			err := submitter.Submit(ctx, Job{
				Handler: "handler",
				Args: Args{
					"seq": i,
				},
			})
			assert.NoError(t, err)
		}
	}()
	go func() {
		defer wg.Done()
		var sum float64
		for i := 0; i < nr; i++ {
			msg, err := subscription.Receive(ctx)
			require.NoError(t, err)
			var job Job
			err = json.Unmarshal(msg.Body, &job)
			require.NoError(t, err)
			assert.Equal(t, "handler", job.Handler)
			seq, ok := job.Args["seq"]
			require.True(t, ok)
			sum += seq.(float64)
			msg.Ack()
		}
		// arithmetic progression
		assert.EqualValues(t, (nr-1)*nr/2, sum)
	}()
	wg.Wait()
}

func TestTopicBadJob(t *testing.T) {
	submitter := NewSubmitter(mempubsub.NewTopic())
	job := Job{
		Handler: "handler",
		// unserializable arg
		Args: Args{"fn": func() {}},
	}
	err := submitter.Submit(context.Background(), job)
	assert.Error(t, err)
}
