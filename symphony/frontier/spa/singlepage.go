// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spa

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/facebookincubator/symphony/pkg/ctxutil"
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/justinas/nosurf"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"go.opencensus.io/plugin/ochttp"
	"go.uber.org/zap"
	"gocloud.dev/gcerrors"
	"gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/httpvar"
	"golang.org/x/xerrors"
)

type (
	// SinglePage is a handler for the single page application.
	SinglePage struct {
		mainjs   string
		manifest *runtimevar.Variable
		template *template.Template
		logger   log.Logger
	}

	// Manifester returns a runtimevar resolving to a manifest.
	Manifester func() (*runtimevar.Variable, error)

	// Option configures single page handler.
	Option func(*SinglePage)

	// index template data.
	indexData struct {
		MainPath   string
		VendorPath string
		Config     []byte
	}

	// renders as json in index template
	configData struct {
		Data appData `json:"appData"`
	}

	// config application data
	appData struct {
		Token    string   `json:"csrfToken"`
		Tabs     []string `json:"tabs"`
		User     []string `json:"user"`
		Features []string `json:"enabledFeatures"`
	}
)

const (
	// distPath is the static directory which content is served
	distPath = "/inventory/static/dist"
	// webpack vendor filename
	vendorjs = "vendor.js"
)

// OriginManifester http gets manifest from origin.
func OriginManifester(origin *url.URL) Manifester {
	endpoint, _ := url.Parse(origin.String())
	endpoint.Path = path.Join(endpoint.Path, distPath, "manifest.json")
	return func() (*runtimevar.Variable, error) {
		return httpvar.OpenVariable(
			&http.Client{
				Transport: &ochttp.Transport{},
				Timeout:   5 * time.Second,
			},
			endpoint.String(),
			runtimevar.NewDecoder(map[string]string{}, runtimevar.JSONDecode),
			nil,
		)
	}
}

// WithLogger sets single page logger.
func WithLogger(logger log.Logger) Option {
	return func(sp *SinglePage) {
		sp.logger = logger
	}
}

//go:generate go run assets_gen.go

// SinglePageHandler creates a handler for the SPA page
func SinglePageHandler(page string, manifester Manifester, opts ...Option) (*SinglePage, error) {
	tmpl, err := vfstemplate.ParseGlob(assets, nil, "*html")
	if err != nil {
		return nil, xerrors.Errorf("loading templates from assets: %w", err)
	}
	manifest, err := manifester()
	if err != nil {
		return nil, xerrors.Errorf("resolving manifest var: %w", err)
	}
	if !strings.HasSuffix(page, ".js") {
		page += ".js"
	}
	sp := &SinglePage{
		mainjs:   page,
		template: tmpl,
		manifest: manifest,
	}
	for _, opt := range opts {
		opt(sp)
	}
	if sp.logger == nil {
		sp.logger = log.NewNopLogger()
	}
	return sp, nil
}

// ServeHTTP implements http.Handler interface.
func (sp *SinglePage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := sp.indexData(r)
	if err != nil {
		sp.logger.For(r.Context()).Error("getting index data", zap.Error(err))
		http.Error(w, "index data unavailable", http.StatusServiceUnavailable)
		return
	}

	if err := sp.template.ExecuteTemplate(w, "index.html", data); err != nil {
		sp.logger.For(r.Context()).Error("executing index template", zap.Error(err))
		http.Error(w, "cannot render index template", http.StatusInternalServerError)
		return
	}
}

// indexData resolves index.html template data from request.
func (sp *SinglePage) indexData(r *http.Request) (*indexData, error) {
	var data indexData
	switch snapshot, err := sp.manifest.Latest(ctxutil.DoneCtx()); {
	case err == nil:
		manifest, ok := snapshot.Value.(map[string]string)
		if !ok {
			return nil, xerrors.Errorf("incompatible manifest type %T", snapshot.Value)
		}
		if data.MainPath, ok = manifest[sp.mainjs]; !ok {
			return nil, xerrors.Errorf("missing manifest key %q", sp.mainjs)
		}
		if data.VendorPath, ok = manifest[vendorjs]; !ok {
			return nil, xerrors.Errorf("missing manifest key %q", vendorjs)
		}
	case gcerrors.Code(err) == gcerrors.NotFound:
		data.MainPath = path.Join(distPath, sp.mainjs)
		data.VendorPath = path.Join(distPath, vendorjs)
	default:
		return nil, xerrors.Errorf("getting latest manifest: %w", err)
	}
	cfg, err := json.Marshal(configData{
		Data: appData{
			Token:    nosurf.Token(r),
			Tabs:     []string{"nms", "inventory"}, // TODO: organization.tabs
			User:     nil,                          // TODO: {"tenant", "email", "isSuperUser"}
			Features: []string{},                   // TODO: feature flags
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("marshaling config data: %w", err)
	}
	data.Config = cfg
	return &data, nil
}

// CheckHealth implements health.Checker interface.
func (sp *SinglePage) CheckHealth() error {
	switch err := sp.manifest.CheckHealth(); gcerrors.Code(err) {
	case gcerrors.OK, gcerrors.NotFound:
		return nil
	default:
		return err
	}
}

// Close implements io.Closer interface.
func (sp *SinglePage) Close() error {
	return sp.manifest.Close()
}
