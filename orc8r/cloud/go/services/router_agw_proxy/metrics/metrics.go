package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
    DeviceLatitude = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_latitude",
        Help: "Device latitude",
    }, []string{"name", "type", "imei", "networkID"})

    DeviceLongitude = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_longitude",
        Help: "Device longitude",
    }, []string{"name", "type", "imei", "networkID"})

    DevicePrivate5GActive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_private5g_active",
        Help: "Private 5G active (1=active, 0=inactive)",
    }, []string{"name", "type", "imei", "networkID"})

    DeviceNetworkService = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_network_service",
        Help: "Network Service (2=5g, 1=4g, 0=unknown)",
    }, []string{"name", "type", "imei", "networkID"})

    DeviceRsrp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_rsrp",
        Help: "Device rsrp value",
    }, []string{"name", "type", "imei", "networkID"})

    DeviceRsrq = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_rsrq",
        Help: "Device rsrq value",
    }, []string{"name", "type", "imei", "networkID"})

    DeviceRssi = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_rssi",
        Help: "Device rssi value",
    }, []string{"name", "type", "imei", "networkID"})

    DeviceSnr = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_snr",
        Help: "Device snr value",
    }, []string{"name", "type", "imei", "networkID"})

    DeviceSpeed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_speed",
        Help: "Device speed value",
    }, []string{"name", "type", "imei", "networkID"})

    DeviceRadioModuleTemp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_radio_module_temp",
        Help: "Device radio module temp",
    }, []string{"name", "type", "imei", "networkID"})

    DeviceBoardTemp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_board_temp",
        Help: "Device board temp",
    }, []string{"name", "type", "imei", "networkID"})

    DeviceCommStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_comm_status",
        Help: "Device Comm Status",
    }, []string{"name", "type", "imei", "networkID"})

    DeviceBytesSent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_bytes_sent",
        Help: "Device bytes sent",
    }, []string{"name", "type", "imei", "networkID"})

    DeviceBytesReceived = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "device_bytes_received",
        Help: "Device bytes received",
    }, []string{"name", "type", "imei", "networkID"})
)

func Init() {
    prometheus.MustRegister(DeviceLatitude, DeviceLongitude, DevicePrivate5GActive, DeviceNetworkService, DeviceRsrp, DeviceRsrq, DeviceRssi, DeviceSnr, DeviceSpeed, DeviceRadioModuleTemp, DeviceBoardTemp, DeviceCommStatus, DeviceBytesSent, DeviceBytesReceived)
}
