// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sentry

import (
	"sync"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/golang/mock/gomock"
	"github.com/magma/magma/src/go/crash"
	"github.com/magma/magma/src/go/crash/sentry/mock_sentry"
	"github.com/stretchr/testify/assert"
)

func TestNewCrash(t *testing.T) {
	options := sentry.ClientOptions{
		Dsn: "http://whatever@really.com/1337",
	}
	c := NewCrash(options)
	crash, ok := c.(*Crash)
	assert.True(
		t, ok, "NewCrash must return *sentry.Crash, have=%+v", c)
	assert.NotNil(t, crash.sentryHub)
	hub, ok := crash.sentryHub.(*sentry.Hub)
	assert.True(
		t, ok, "crash.sentryHub must contain a sentry.hub, have=%+v", crash.sentryHub)

	assert.Equal(t, options.Dsn, hub.Client().Options().Dsn)
}

func TestNewCrash_Error(t *testing.T) {
	options := sentry.ClientOptions{
		Dsn: "\"bad_input",
	}
	assert.Panics(t, func() {
		_ = NewCrash(options)
	})
}

func TestCrash_AddBreadcrumb(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedType := "test type"
	expectedCategory := "test category"
	expectedMsg := "test msg"
	expectedLevel := sentry.LevelInfo
	expectedTimestamp := time.Now()

	expectedBreadcrumb := &sentry.Breadcrumb{
		Type:      expectedType,
		Category:  expectedCategory,
		Message:   expectedMsg,
		Level:     expectedLevel,
		Timestamp: expectedTimestamp,
	}

	mockHub := mock_sentry.NewMocksentryHub(ctrl)
	mockHub.EXPECT().AddBreadcrumb(expectedBreadcrumb, nil)

	c := &Crash{
		sentryHub: mockHub,
	}

	c.AddBreadcrumb(crash.Breadcrumb{
		Type:      expectedType,
		Category:  expectedCategory,
		Message:   expectedMsg,
		Level:     crash.Level(expectedLevel),
		Timestamp: expectedTimestamp,
	})
}

func TestCrash_AddBreadcrumb_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	breadcrumb := crash.Breadcrumb{}

	mockHub := mock_sentry.NewMocksentryHub(ctrl)
	mockHub.EXPECT().AddBreadcrumb(&sentry.Breadcrumb{}, nil)

	c := &Crash{
		sentryHub: mockHub,
	}

	c.AddBreadcrumb(breadcrumb)
}

func TestCrash_Recover_NoError(t *testing.T) {
	c := &Crash{}
	assert.Equal(t, crash.EventID(""), c.Recover(nil))
}

func TestCrash_Recover(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHub := mock_sentry.NewMocksentryHub(ctrl)

	var id sentry.EventID = "testEventID"
	mockHub.EXPECT().Recover("trigger").Return(&id)

	c := &Crash{
		sentryHub: mockHub,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			err := recover()
			assert.NotNil(t, err)
			eventID := c.Recover(err)
			assert.Equal(t, crash.EventID(id), eventID)
		}()
		panic("trigger")
	}()
	wg.Wait()
}

func TestCrash_Flush_TimeoutExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeout := time.Second

	mockHub := mock_sentry.NewMocksentryHub(ctrl)
	mockHub.EXPECT().Flush(timeout)

	c := &Crash{
		sentryHub: mockHub,
	}

	assert.False(t, c.Flush(timeout))
}

func TestCrash_Flush_TimeoutNotExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeout := time.Minute

	mockHub := mock_sentry.NewMocksentryHub(ctrl)

	mockHub.EXPECT().Flush(timeout).Return(true)

	c := &Crash{
		sentryHub: mockHub,
	}

	assert.True(t, c.Flush(timeout))
}

func TestSentry_convertSentryEventID_Nil(t *testing.T) {
	assert.Equal(t, crash.EventID(""), convertSentryEventID(nil))
}

func TestSentry_convertSentryEventID(t *testing.T) {
	eventIDStr := "test"
	sentryEventID := sentry.EventID(eventIDStr)
	assert.Equal(t, crash.EventID(eventIDStr), convertSentryEventID(&sentryEventID))
}
