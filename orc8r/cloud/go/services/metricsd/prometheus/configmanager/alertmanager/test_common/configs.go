package test_common

import (
	"net/url"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/config"

	amconfig "github.com/prometheus/alertmanager/config"
)

var (
	sampleURL, _ = url.Parse("http://test.com")
	SampleRoute  = config.Route{
		Receiver: "testReceiver",
		Routes: []*config.Route{
			{
				Receiver: "testReceiver",
			},
			{
				Receiver: "slack_receiver",
			},
		},
	}
	SampleReceiver = config.Receiver{
		Name: "testReceiver",
	}
	SampleSlackReceiver = config.Receiver{
		Name: "slack_receiver",
		SlackConfigs: []*config.SlackConfig{{
			APIURL:   "http://slack.com/12345",
			Username: "slack_user",
			Channel:  "slack_alert_channel",
		}},
	}
	SamplePagerDutyReceiver = config.Receiver{
		Name: "pagerduty_receiver",
		PagerDutyConfigs: []*config.PagerDutyConfig{{
			ServiceKey: "0",
		}},
	}
	SamplePushoverReceiver = config.Receiver{
		Name: "pushover_receiver",
		PushoverConfigs: []*config.PushoverConfig{{
			UserKey: "101",
			Token:   "1",
		}},
	}
	SampleWebhookReceiver = config.Receiver{
		Name: "webhook_receiver",
		WebhookConfigs: []*config.WebhookConfig{{
			URL: &amconfig.URL{
				URL: sampleURL,
			},
			NotifierConfig: amconfig.NotifierConfig{
				VSendResolved: true,
			},
		}},
	}
	SampleEmailReceiver = config.Receiver{
		Name: "email_receiver",
		EmailConfigs: []*config.EmailConfig{{
			To:        "test@mail.com",
			From:      "sampleUser",
			Headers:   map[string]string{"header": "value"},
			Smarthost: "http://mail-server.com",
		}},
	}
	SampleConfig = config.Config{
		Route: &SampleRoute,
		Receivers: []*config.Receiver{
			&SampleSlackReceiver, &SampleReceiver, &SamplePagerDutyReceiver, &SamplePushoverReceiver, &SampleWebhookReceiver, &SampleEmailReceiver,
		},
	}
)
