package types

type APIResponse struct {
    Items []DeviceItem `json:"items"`
}

type DeviceItem struct {
    Data       DeviceData `json:"data"`
    Name       string     `json:"name"`
    Type       string     `json:"type"`
    CommStatus string     `json:"commStatus"`
}

type DeviceData struct {
    Latitude           *float64 `json:"latitude,string"`
    Longitude          *float64 `json:"longitude,string"`
    CellId             string   `json:"cellId"`
    Imei               string   `json:"imei"`
    SignalStrength     string   `json:"signalStrength"`
    BytesSent          *string  `json:"bytesSent"`
    BytesReceived      *string  `json:"bytesReceived"`
    NetworkServiceType string   `json:"networkServiceType"`
    Rsrp               *string  `json:"rsrp"`
    Rsrq               *string  `json:"rsrq"`
    Rssi               *string  `json:"rssi"`
    Snr                *string  `json:"snr"`
    Speed              *string  `json:"speed"`
    RadioModuleTemp    *string  `json:"radioModuleTemp"`
    BoardTemp          *string  `json:"boardTemp"`
}
