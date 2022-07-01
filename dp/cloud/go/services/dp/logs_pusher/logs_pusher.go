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
}

func PushDPLog(ctx context.Context, log *DPLog, consumerUrl string) error {
	client := http.Client{}
	body, err := json.Marshal(*log)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest(http.MethodPost, consumerUrl, strings.NewReader(string(body)))
	req.Header.Set("cntentType", "application/json")
	req.WithContext(ctx)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
