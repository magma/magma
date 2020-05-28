// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/dialect"
)

// Option function to configure the client.
type Option func(*config)

// Config is the configuration for the client and its builder.
type config struct {
	// driver used for executing database requests.
	driver dialect.Driver
	// debug enable a debug logging.
	debug bool
	// log used for logging on debug mode.
	log func(...interface{})
	// hooks to execute on mutations.
	hooks *hooks
}

// hooks per client, for fast access.
type hooks struct {
	ActionsRule                 []ent.Hook
	Activity                    []ent.Hook
	CheckListCategory           []ent.Hook
	CheckListCategoryDefinition []ent.Hook
	CheckListItem               []ent.Hook
	CheckListItemDefinition     []ent.Hook
	Comment                     []ent.Hook
	Customer                    []ent.Hook
	Equipment                   []ent.Hook
	EquipmentCategory           []ent.Hook
	EquipmentPort               []ent.Hook
	EquipmentPortDefinition     []ent.Hook
	EquipmentPortType           []ent.Hook
	EquipmentPosition           []ent.Hook
	EquipmentPositionDefinition []ent.Hook
	EquipmentType               []ent.Hook
	File                        []ent.Hook
	FloorPlan                   []ent.Hook
	FloorPlanReferencePoint     []ent.Hook
	FloorPlanScale              []ent.Hook
	Hyperlink                   []ent.Hook
	Link                        []ent.Hook
	Location                    []ent.Hook
	LocationType                []ent.Hook
	PermissionsPolicy           []ent.Hook
	Project                     []ent.Hook
	ProjectType                 []ent.Hook
	Property                    []ent.Hook
	PropertyType                []ent.Hook
	ReportFilter                []ent.Hook
	Service                     []ent.Hook
	ServiceEndpoint             []ent.Hook
	ServiceEndpointDefinition   []ent.Hook
	ServiceType                 []ent.Hook
	Survey                      []ent.Hook
	SurveyCellScan              []ent.Hook
	SurveyQuestion              []ent.Hook
	SurveyTemplateCategory      []ent.Hook
	SurveyTemplateQuestion      []ent.Hook
	SurveyWiFiScan              []ent.Hook
	User                        []ent.Hook
	UsersGroup                  []ent.Hook
	WorkOrder                   []ent.Hook
	WorkOrderDefinition         []ent.Hook
	WorkOrderType               []ent.Hook
}

// Options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...interface{})) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}
