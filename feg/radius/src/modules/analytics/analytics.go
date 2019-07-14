/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package analytics

import (
	"crypto/sha1"
	"crypto/tls"
	"fbc/cwf/radius/session"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"

	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fbc/lib/go/radius/rfc2866"
	"fbc/lib/go/radius/rfc2869"

	"fbc/lib/go/libgraphql"

	"go.uber.org/zap"
)

// Config configuration structure for the EAP module
type Config struct {
	// the Access Token for WWW GraphQL calls
	AccessToken string
	// the URL for WWW GraphQL calls
	GraphQLURL string
	// true means all GraphQL operations will be eliminated & assumed successful.
	DryRunGraphQL bool
}

// the name of this module for logging context
const moduleName = "analytics"

var (
	// a client to WWW GraphQL with which were calling WWW Ops
	graphqlClient *libgraphql.Client
	// true means all GraphQL operations will be eliminated & assumed successful.
	cfg Config
)

// set the defaults for all configuration parameters
func setDefaultCfg() {
	cfg = Config{
		GraphQLURL:    "https://graph.expresswifi.com/graphql",
		DryRunGraphQL: false,
	}
}

// Init module interface implementation
//nolint:deadcode
func Init(logger *zap.Logger, config modules.ModuleConfig) error {
	setDefaultCfg()
	// Parse config
	err := mapstructure.Decode(config, &cfg)
	if err != nil {
		return err
	}
	if cfg.DryRunGraphQL {
		logger.Warn("GraphQL set as dry-run")
	}
	graphqlClient = libgraphql.NewClient(libgraphql.ClientConfig{
		Token:    cfg.AccessToken,
		Endpoint: cfg.GraphQLURL,
		HTTPClient: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}},
	})

	return nil
}

// get the session state when exist. otherwise, create it
func getSessionState(logger *zap.Logger, c *modules.RequestContext) (*session.State, error) {
	sessionState, err := c.SessionStorage.Get()
	if err == nil {
		// session state found
		return sessionState, nil
	}
	// when we fail to access the DB, i prefer loosing an update bcz most will be AccountingUpdate rather than Auth.
	// the next AccountingUpdate will catch-up the missing one (assuming transient errors & recovery of services)
	// creating a new session in case of error means a DB failure can cause a storm of new sessions. plus, not all info is available in Accounting vs. Auth
	if err == session.ErrInvalidDataFormat {
		return nil, err
	}
	logger.Debug("creating default session state")
	sessionState = &session.State{RadiusSessionFBID: 0}
	return sessionState, nil
}

// do the GraphQL call to create the RadiusSession
func createRadiusSession(logger *zap.Logger, c *modules.RequestContext, session *RadiusSession, sessionState *session.State) {
	logger.Debug("Creating a new RADIUS session", zap.Any("radius_session", session))
	if cfg.DryRunGraphQL {
		sessionState.RadiusSessionFBID = uint64(time.Now().UnixNano())
		time.Sleep(time.Millisecond) // provide some delay for GraphQL calls
		logger.Debug("GraphQL is in dry-run mode !!!", zap.Any("radius_session", session))
	} else {
		createOp := NewCreateSessionOp(session)
		err := graphqlClient.Do(createOp)
		if err != nil {
			logger.Error("failed creating session", zap.Any("radius_session", &session),
				zap.Error(err))
			return
		}
		sessionState.RadiusSessionFBID = createOp.Response().FBID
	}
	err := c.SessionStorage.Set(*sessionState)
	if err != nil {
		logger.Error("failed to update session", zap.Error(err), zap.Any("session_state", sessionState))
	}
}

// do the GraphQL call to update the RadiusSession
func updateRadiusSession(logger *zap.Logger, c *modules.RequestContext, session *RadiusSession, sessionState *session.State) {
	if cfg.DryRunGraphQL {
		time.Sleep(time.Millisecond) // provide some delay for GraphQL calls
		logger.Debug("GraphQL is in dry-run mode !!!", zap.Any("radius_session", session))
	} else {
		updateOp := NewUpdateSessionOp(session)
		err := graphqlClient.Do(updateOp)
		if err != nil {
			logger.Error("failed updating session", zap.Any("radius_session", &session),
				zap.Error(err))
			return
		}
	}
	err := c.SessionStorage.Set(*sessionState) // TODO: do we need to save the accounting state every time its updated ???
	if err != nil {
		logger.Error("failed to update session", zap.Error(err), zap.Any("session_state", sessionState))
	}
}

// Handle module interface implementation
//nolint:deadcode
func Handle(c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	pkt := r.Packet
	logger := c.Logger.With(zap.String("module", moduleName), zap.Int("radius_code", int(r.Code)))
	var asyncQL sync.WaitGroup

	switch r.Code {
	case radius.CodeAccessRequest:
		sessionState, err := getSessionState(logger, c)
		if err != nil {
			logger.Error("failed to get session state", zap.Any("radius_request", r), zap.Error(err))
			break
		}
		framedIPAddr := fmt.Sprintf("%v", rfc2865.FramedIPAddress_Get(pkt))
		nasIDAddr := fmt.Sprintf("%v", rfc2865.NASIPAddress_Get(pkt))
		calledStationID := rfc2865.CalledStationID_GetString(pkt)
		calledStationIDSeparator := strings.IndexByte(calledStationID, ':')
		normalizedMacAddress := calledStationID
		if calledStationIDSeparator != -1 {
			// remove trailing ":<AP name>", format is "AB-CD-EF-GH-IJ-KL", I.E.: MAC address, in capitals
			normalizedMacAddress = strings.ToUpper(normalizedMacAddress[:calledStationIDSeparator])
		}
		logger.Info("processing auth packet", zap.String("framed_ip_addr", framedIPAddr),
			zap.String("nas_ip_addr", nasIDAddr))
		session := RadiusSession{
			NASIPAddress:         nasIDAddr,
			NASIdentifier:        tokenizeString(rfc2865.NASIdentifier_GetString(pkt)),
			AcctSessionID:        tokenizeString(rfc2866.AcctSessionID_GetString(pkt)),
			CalledStationID:      tokenizeString(calledStationID),
			CallingStationID:     tokenizeString(rfc2865.CallingStationID_GetString(pkt)),
			FramedIPAddress:      tokenizeString(framedIPAddr),
			NormalizedMacAddress: tokenizeString(normalizedMacAddress),
		}
		asyncQL.Add(1)
		go func() {
			createRadiusSession(logger, c, &session, sessionState)
			asyncQL.Done()
		}()

	case radius.CodeAccountingRequest:
		switch rfc2866.AcctStatusType_Get(pkt) {
		case rfc2866.AcctStatusType_Value_Start:
			fallthrough
		case rfc2866.AcctStatusType_Value_Stop:
			fallthrough
		case rfc2866.AcctStatusType_Value_InterimUpdate:
			sessionState, err := c.SessionStorage.Get()
			if err != nil {
				logger.Error("failed to get session state", zap.Any("radius_request", r), zap.Error(err))
				break
			}
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
			logger.Info("processing accounting packet", zap.Int64("input_bytes", inputBytes),
				zap.Int64("output_bytes", outputBytes))
			session := RadiusSession{
				FBID:          sessionState.RadiusSessionFBID,
				NASIdentifier: rfc2865.NASIdentifier_GetString(pkt),
				AcctSessionID: rfc2866.AcctSessionID_GetString(pkt),
				UploadBytes:   inputBytes,
				DownloadBytes: outputBytes,
			}
			asyncQL.Add(1)
			go func() {
				updateRadiusSession(logger, c, &session, sessionState)
				asyncQL.Done()
			}()

		case rfc2866.AcctStatusType_Value_AccountingOn:
			fallthrough
		case rfc2866.AcctStatusType_Value_AccountingOff:
			fallthrough
		case rfc2866.AcctStatusType_Value_Failed:
			rfc2866.AcctStatusType_Get(pkt)
			logger.Info("ignoring accounting packet",
				zap.String("acct_status_type", rfc2866.AcctStatusType_Get(pkt).String()))
		}
	}
	// since we only provide analytics, dont fail packet processing even when we have errors, so client flow isn't hampered.
	resp, err := next(c, r)
	// wait for the GraphQL we fired before the module chain was called
	asyncQL.Wait()
	return resp, err
}

// tokenizeString tokenize a PII string field
func tokenizeString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return string(h.Sum(nil))
}
