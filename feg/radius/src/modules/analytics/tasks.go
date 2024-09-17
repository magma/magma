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

package analytics

import (
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/session"
	"time"

	"go.uber.org/zap"
)

type (

	// a task to wrap a GraphQL "CreateSession" operation to be serialized with all operations of this session
	createSessionTask struct {
		Logger             *zap.Logger
		reqCtx             *modules.RequestContext
		session            *RadiusSession
		untokenizedSession *RadiusSession
		sessionState       *session.State
	}

	// a task to wrap a GraphQL "UpdateSession" operation to be serialized with all operations of this session
	updateSessionTask struct {
		Logger       *zap.Logger
		reqCtx       *modules.RequestContext
		session      *RadiusSession
		sessionState *session.State
	}
)

// Run to be called by task execution engine
func (op *createSessionTask) Run(m ModuleCtx) {
	createRadiusSession(m, op.Logger, op.reqCtx, op.session, op.sessionState)
}

// Run to be called by task execution engine
func (op *updateSessionTask) Run(m ModuleCtx) {
	updateRadiusSession(m, op.Logger, op.reqCtx, op.session, op.sessionState)
}

// do the GraphQL call to create the RadiusSession
func createRadiusSession(m ModuleCtx, logger *zap.Logger, c *modules.RequestContext, session *RadiusSession, sessionState *session.State) {
	logger.Debug("Creating a new RADIUS session", zap.Any("radius_session", session))
	if m.cfg.DryRunGraphQL {
		sessionState.RadiusSessionFBID = uint64(time.Now().UnixNano())
		time.Sleep(time.Millisecond) // provide some delay for GraphQL calls
		logger.Debug("GraphQL is in dry-run mode !!!", zap.Any("radius_session", session))
	} else {
		createOp := NewCreateSessionOp(session)
		err := m.graphqlClient.Do(createOp)
		if err != nil {
			logger.Error("failed creating session", zap.Any("radius_session", &session),
				zap.Error(err))
			return
		}
		logger.Warn("session created", zap.Uint64("fbid", sessionState.RadiusSessionFBID))
		sessionState.RadiusSessionFBID = createOp.Response().FBID
	}

	// Persist state
	err := c.SessionStorage.Set(*sessionState)
	if err != nil {
		logger.Error("failed to update session", zap.Error(err), zap.Any("session_state", sessionState))
	}
}

// do the GraphQL call to update the RadiusSession
func updateRadiusSession(m ModuleCtx, logger *zap.Logger, c *modules.RequestContext, session *RadiusSession, sessionState *session.State) {
	logger.Debug("Updating RADIUS session", zap.Any("radius_session", session))
	if m.cfg.DryRunGraphQL {
		time.Sleep(time.Millisecond) // provide some delay for GraphQL calls
		logger.Debug("GraphQL is in dry-run mode !!!", zap.Any("radius_session", session))
	} else {
		updateOp := NewUpdateSessionOp(session)
		err := m.graphqlClient.Do(updateOp)
		if err != nil {
			logger.Error("failed updating session", zap.Any("radius_session", &session),
				zap.Error(err))
			return
		}
	}

	// Persist state
	err := c.SessionStorage.Set(*sessionState)
	if err != nil {
		logger.Error("failed to update session", zap.Error(err), zap.Any("session_state", sessionState))
	}
}

// push a lazy task to the queue
// create the queue if it doesnt exist yet
func pushGraphQLTask(m ModuleCtx, logger *zap.Logger, task Request, sessionState *session.State) {
	// create a queue if not exist
	q := m.graphQLOps[sessionState.AcctSessionID]
	if q == nil {
		// no queue for the session - create one
		q = NewAnalyticsQueue(m)
		m.graphQLOps[sessionState.AcctSessionID] = q
		logger.Debug("Creating graphql queue", zap.String("acct_session_id", sessionState.AcctSessionID))
	}
	q.Push(task)
}

// cleanSessionTasks drain tasks & delete the queue from the map
func cleanSessionTasks(m ModuleCtx, logger *zap.Logger, sessionState *session.State) {
	q := m.graphQLOps[sessionState.AcctSessionID]
	if q == nil {
		return
	}
	q.Close(true)
	delete(m.graphQLOps, sessionState.AcctSessionID)
	logger.Debug("drained graphql queue", zap.String("acct_session_id", sessionState.AcctSessionID),
		zap.Bool("is_exist", q != nil))
}
