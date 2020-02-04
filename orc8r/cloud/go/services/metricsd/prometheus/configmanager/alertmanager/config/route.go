package config

import (
	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/model"
)

// Route provides a struct to marshal/unmarshal into an alertmanager route
// since that struct does not support json encoding
type Route struct {
	Receiver string `yaml:"receiver,omitempty" json:"receiver,omitempty"`

	GroupByStr []string          `yaml:"group_by,omitempty" json:"group_by,omitempty"`
	GroupBy    []model.LabelName `yaml:"-" json:"-"`
	GroupByAll bool              `yaml:"-" json:"-"`

	Match    map[string]string        `yaml:"match,omitempty" json:"match,omitempty"`
	MatchRE  map[string]config.Regexp `yaml:"match_re,omitempty" json:"match_re,omitempty"`
	Continue bool                     `yaml:"continue,omitempty" json:"continue,omitempty"`
	Routes   []*Route                 `yaml:"routes,omitempty" json:"routes,omitempty"`

	GroupWait      string `yaml:"group_wait,omitempty" json:"group_wait,omitempty"`
	GroupInterval  string `yaml:"group_interval,omitempty" json:"group_interval,omitempty"`
	RepeatInterval string `yaml:"repeat_interval,omitempty" json:"repeat_interval,omitempty"`
}
