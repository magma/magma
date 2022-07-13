package logs_pusher

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type DPLog struct {
	EventTimestamp   int64  `json:"event_timestamp"`
	LogFrom          string `json:"log_from"`
	LogTo            string `json:"log_to"`
	LogName          string `json:"log_name"`
	LogMessage       string `json:"log_message"`
	CbsdSerialNumber string `json:"cbsd_serial_number"`
	NetworkId        string `json:"network_id"`
	FccId            string `json:"fcc_id"`
}

type LogPusher func(ctx context.Context, log *DPLog, consumerUrl string) error

func PushDPLog(ctx context.Context, log *DPLog, consumerUrl string) error {
	client := http.Client{}
	body, _ := json.Marshal(log)
	req, _ := http.NewRequest(http.MethodPost, consumerUrl, strings.NewReader(string(body)))
	req.Header.Set("contentType", "application/json")
	_, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	return nil
}
