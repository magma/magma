// +build !go1.8

package diam

import "crypto/tls"

func TLSConfigClone(cfg *tls.Config) *tls.Config {
	if cfg != nil {
		newCfg := *cfg
		return &newCfg
	}
	return nil
}
