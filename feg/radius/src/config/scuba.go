package config

// Scuba scuba logger config
type Scuba struct {
	MessageQueueSize int    `json:"message_queue_size" default:"2000"`
	FlushIntervalSec int    `json:"flush_interval_sec" default:"2"`
	BatchSize        int    `json:"batch_size" default:"15"`
	GraphURL         string `json:"graph_url" default:"https://graph.facebook.com/scribe_logs"`
	AccessToken      string `json:"access_token"`
	PartnerShortName string `json:"partner_shortname"`
}
