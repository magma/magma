package config

// Ods confoguration for ODS bindings
type Ods struct {
	Category        string   `json:"category_id" required:"true"`
	Prefix          string   `json:"prefix" required:"true"`
	Token           string   `json:"access_token" required:"true"`
	Entity          string   `json:"entity" required:"true"`
	DisablePost     bool     `json:"disable_post" default:"false"`
	DebugPrints     bool     `json:"debug_prints" default:"false"`
	GraphURL        string   `json:"graph_url" default:"https://graph.facebook.com/ods_metrics"`
	ReportingPeriod Duration `json:"reporting_period" default:"60s"`
}
