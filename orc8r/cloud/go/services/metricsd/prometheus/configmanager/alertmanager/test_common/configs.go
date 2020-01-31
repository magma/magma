package test_common

import (
	"net/url"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/config"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/receivers"

	amconfig "github.com/prometheus/alertmanager/config"
)

var (
	sampleURL, _ = url.Parse("http://test.com")
	SampleRoute  = amconfig.Route{
		Receiver: "testReceiver",
		Routes: []*amconfig.Route{
			{
				Receiver: "testReceiver",
			},
			{
				Receiver: "slack_receiver",
			},
		},
	}
	SampleReceiver = receivers.Receiver{
		Name: "testReceiver",
	}
	SampleSlackReceiver = receivers.Receiver{
		Name: "slack_receiver",
		SlackConfigs: []*receivers.SlackConfig{{
			APIURL:   "http://slack.com/12345",
			Username: "slack_user",
			Channel:  "slack_alert_channel",
		}},
	}
	SampleWebhookReceiver = receivers.Receiver{
		Name: "webhook_receiver",
		WebhookConfigs: []*receivers.WebhookConfig{{
			URL: &amconfig.URL{
				URL: sampleURL,
			},
			NotifierConfig: amconfig.NotifierConfig{
				VSendResolved: true,
			},
		}},
	}
	SampleEmailReceiver = receivers.Receiver{
		Name: "email_receiver",
		EmailConfigs: []*receivers.EmailConfig{{
			To:        "test@mail.com",
			From:      "sampleUser",
			Headers:   map[string]string{"header": "value"},
			Smarthost: "http://mail-server.com",
		}},
	}
	SampleConfig = config.Config{
		Route: &SampleRoute,
		Receivers: []*receivers.Receiver{
			&SampleSlackReceiver, &SampleReceiver, &SampleWebhookReceiver, &SampleEmailReceiver,
		},
	}
)
