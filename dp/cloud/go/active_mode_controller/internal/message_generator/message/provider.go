package message

import (
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/active_mode_controller/protos/requests"
)

type ClientProvider interface {
	GetRequestsClient() requests.RadioControllerClient
	GetActiveModeClient() active_mode.ActiveModeControllerClient
}
