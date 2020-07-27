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
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/modules/analytics/graphql"
	"fbc/cwf/radius/session"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fbc/lib/go/radius/rfc2866"
	"fbc/lib/go/radius/rfc2869"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

// Config configuration structure for the EAP module
type Config struct {
	AccessToken   string // Access Token to use in GraphQL calls
	GraphQLURL    string // the GraphQL endpoint to issue calls to
	DryRunGraphQL bool   // true means all GraphQL operations will be skipped & assumed successful.
	AllowPII      bool   // If true, PII will not be tokenized before sending to GraphQL
}

type (
	// ModuleCtx ...
	ModuleCtx struct {
		// a client to issue GraphQL calls
		graphqlClient *graphql.Client
		cfg           Config
		graphQLOps    map[string]*Queue // a queue of GraphQL operations, pending serialized execution
	}
)

// Init module interface implementation
//nolint:deadcode
func Init(logger *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	var ctx ModuleCtx

	// Parse config
	err := mapstructure.Decode(config, &ctx.cfg)
	if err != nil {
		return nil, err
	}

	// Warn in log about dangerous settings
	if ctx.cfg.DryRunGraphQL {
		logger.Warn("ANALYTICS IS SET TO DRY MODE, DATA WILL NOT BE SENT OUT VIA GRAPHQL")
	}

	if ctx.cfg.AllowPII {
		logger.Warn("ANALYTICS IS SET TO ALLOW PII BE SENT OUT")
	}

	ctx.graphQLOps = make(map[string]*Queue)
	// Create client
	ctx.graphqlClient = graphql.NewClient(graphql.ClientConfig{
		Token:    ctx.cfg.AccessToken,
		Endpoint: ctx.cfg.GraphQLURL,
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	})

	return ctx, nil
}

// get the session state when exist. otherwise, create it
func getSessionState(logger *zap.Logger, c *modules.RequestContext) (*session.State, error) {
	sessionState, err := c.SessionStorage.Get()
	if err == nil {
		return sessionState, nil
	}

	// When we fail to access the DB, we prefer loosing an update bcz most will be
	// AccountingUpdate rather than Auth.
	// the next AccountingUpdate will catch-up the missing one (assuming transient errors
	// & recovery of services) creating a new session in case of error means a DB failure
	// can cause a storm of new sessions. plus, not all info is available in Acct vs. Auth
	if err == session.ErrInvalidDataFormat {
		return nil, err
	}

	logger.Debug("creating default session state")
	sessionState = &session.State{RadiusSessionFBID: 0}
	return sessionState, nil
}

// Handle module interface implementation
//nolint:deadcode
func Handle(m modules.Context, c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	mCtx := m.(ModuleCtx)
	pkt := r.Packet

	switch r.Code {
	case radius.CodeAccessRequest:
		// Get session state
		sessionState, err := getSessionState(c.Logger, c)
		if err != nil {
			c.Logger.Error("failed to get session state", zap.Any("radius_request", r), zap.Error(err))
			break
		}

		// Build Session structure
		framedIPAddr := fmt.Sprintf("%v", rfc2865.FramedIPAddress_Get(pkt))
		nasIDAddr := fmt.Sprintf("%v", rfc2865.NASIPAddress_Get(pkt))
		calledStationID := rfc2865.CalledStationID_GetString(pkt)
		calledStationIDSeparator := strings.IndexByte(calledStationID, ':')
		normalizedMacAddress := calledStationID
		if calledStationIDSeparator != -1 {
			// remove trailing ":<AP name>", format is "AB-CD-EF-GH-IJ-KL", I.E.: MAC address, in capitals
			normalizedMacAddress = strings.ToUpper(normalizedMacAddress[:calledStationIDSeparator])
		}
		c.Logger.Info("processing auth packet", zap.String("framed_ip_addr", framedIPAddr),
			zap.String("nas_ip_addr", nasIDAddr))
		session := RadiusSession{
			NASIPAddress:         nasIDAddr,
			NASIdentifier:        rfc2865.NASIdentifier_GetString(pkt),
			AcctSessionID:        rfc2866.AcctSessionID_GetString(pkt),
			CalledStationID:      calledStationID,
			CallingStationID:     rfc2865.CallingStationID_GetString(pkt),
			FramedIPAddress:      framedIPAddr,
			NormalizedMacAddress: normalizedMacAddress,
		}

		if !mCtx.cfg.AllowPII {
			session.AcctSessionID = tokenize(session.AcctSessionID)
			session.CallingStationID = tokenize(session.CallingStationID)
			session.FramedIPAddress = tokenize(session.FramedIPAddress)
			session.NormalizedMacAddress = tokenize(session.NormalizedMacAddress)
		}

		// Persist state before we fire an async task - so the session has a single state object
		err = c.SessionStorage.Set(*sessionState)
		if err != nil {
			c.Logger.Error("failed to update session", zap.Error(err), zap.Any("session_state", sessionState))
		}

		pushGraphQLTask(mCtx, c.Logger, &createSessionTask{
			Logger:       c.Logger,
			reqCtx:       c,
			session:      &session,
			sessionState: sessionState,
		}, sessionState)

	case radius.CodeAccountingRequest:
		switch rfc2866.AcctStatusType_Get(pkt) {
		case rfc2866.AcctStatusType_Value_Start:
			fallthrough
		case rfc2866.AcctStatusType_Value_Stop:
			fallthrough
		case rfc2866.AcctStatusType_Value_InterimUpdate:
			sessionState, err := c.SessionStorage.Get()
			if err != nil {
				c.Logger.Error("failed to get session state", zap.Any("radius_request", r), zap.Error(err))
				break
			}

			// Extract accounting octets
			inputBytes := int64(rfc2866.AcctInputOctets_Get(pkt))
			inputWrapped := rfc2869.AcctInputGigawords_Get(pkt)
			if inputWrapped != 0 {
				inputBytes |= int64(inputWrapped) << 32
			}
			outputBytes := int64(rfc2866.AcctOutputOctets_Get(pkt))
			outputWrapped := rfc2869.AcctOutputGigawords_Get(pkt)
			if outputWrapped != 0 {
				outputBytes |= int64(outputWrapped) << 32
			}
			c.Logger.Debug(
				"processing accounting packet",
				zap.Int64("input_bytes", inputBytes),
				zap.Int64("output_bytes", outputBytes),
			)

			// Extract accounting octets
			session := RadiusSession{
				FBID:          sessionState.RadiusSessionFBID,
				NASIdentifier: rfc2865.NASIdentifier_GetString(pkt),
				AcctSessionID: rfc2866.AcctSessionID_GetString(pkt),
				UploadBytes:   inputBytes,
				DownloadBytes: outputBytes,
			}

			// Tokenize fields which might contain PII
			if !mCtx.cfg.AllowPII {
				session.AcctSessionID = tokenize(session.AcctSessionID)
			}

			// Send the request!
			pushGraphQLTask(mCtx, c.Logger, &updateSessionTask{
				Logger:       c.Logger,
				reqCtx:       c,
				session:      &session,
				sessionState: sessionState,
			}, sessionState)
			if rfc2866.AcctStatusType_Get(pkt) == rfc2866.AcctStatusType_Value_Stop {
				cleanSessionTasks(mCtx, c.Logger, sessionState)
			}
		case rfc2866.AcctStatusType_Value_AccountingOn:
			fallthrough
		case rfc2866.AcctStatusType_Value_AccountingOff:
			fallthrough
		case rfc2866.AcctStatusType_Value_Failed:
			acctStatusType := rfc2866.AcctStatusType_Get(pkt)
			c.Logger.Info(
				"ignoring accounting packet",
				zap.String("acct_status_type", acctStatusType.String()),
			)
		}
	}
	// since we only provide analytics, dont fail packet processing even when we have errors, so client flow isn't hampered.
	resp, err := next(c, r)
	return resp, err
}

// tokenize tokenize a PII string field
func tokenize(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
