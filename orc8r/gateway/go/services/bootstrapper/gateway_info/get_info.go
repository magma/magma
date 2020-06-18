// package gateway_info collects, formats & returns GW information needed for gateway registration
package gateway_info

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"magma/gateway/config"
	"magma/gateway/services/bootstrapper/service"
	"magma/orc8r/lib/go/security/key"

	"github.com/emakeev/snowflake"
)

// Get returns Gateway Hardware Id and bootstrapping public key
func Get() (hwId string, pubKey interface{}, err error) {
	uuid, err := snowflake.Make()
	hwId = uuid.String()
	if err != nil {
		return
	}
	ck, err := service.GetChallengeKey()
	if err != nil {
		return
	}
	pubKey = key.PublicKey(ck)
	return
}

// GetFormatted returns formatted string with GW information
func GetFormatted() (string, error) {
	hwId, pubKey, err := Get()
	if err != nil {
		if len(hwId) > 0 {
			err = fmt.Errorf("failed to get challenge key for GW ID %s: %v", hwId, err)
		}
		return "", err
	}
	res := fmt.Sprintf("\nHardware ID:\n------------\n%s\n", hwId)
	marshaledPubKey, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		err = fmt.Errorf("failed to marshal Public Challenge Key from %s: %v",
			config.GetMagmadConfigs().BootstrapConfig.ChallengeKey, err)
		return res, err
	}
	res = fmt.Sprintf("%s\nChallenge Key:\n--------------\n%s\n",
		res, base64.StdEncoding.EncodeToString(marshaledPubKey))
	return res, nil
}
