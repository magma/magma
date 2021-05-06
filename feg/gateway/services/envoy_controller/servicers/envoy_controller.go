/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package servicers implements the grpc logic for the Envoy Controller
package servicers

import (
	"context"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/envoy_controller/control_plane"

	"github.com/golang/glog"
)

type envoyControllerService struct {
	ueInfos       control_plane.UEInfoMap
	controllerCli control_plane.EnvoyController
}

// AddUEHeaderEnrichment adds the UE to the current header enrichment list, if UE is already in the list replaces the he information for that UE
func (s *envoyControllerService) AddUEHeaderEnrichment(
	ctx context.Context,
	req *protos.AddUEHeaderEnrichmentRequest,
) (*protos.AddUEHeaderEnrichmentResult, error) {
	glog.Infof("AddUEHeaderEnrichmentResult received for IP %s", req.UeIp.Address)
	glog.V(2).Infof("req %s", req)

	ueIp := string(req.UeIp.Address)

	if _, ok := s.ueInfos[ueIp]; !ok {
		s.ueInfos[ueIp] = map[string]*control_plane.UEInfo{}
	} else {
		if _, ok := s.ueInfos[ueIp][req.RuleId]; ok {
			return &protos.AddUEHeaderEnrichmentResult{Result: protos.AddUEHeaderEnrichmentResult_RULE_ID_CONFLICT}, nil
		}
	}

	// Loop over other rules, check if there is a conflict with the Websites (2 identical Websites will cause an envoy deadloop)
	for _, ue_info := range s.ueInfos[ueIp] {
		for _, new_website := range req.Websites {
			for _, existing_website := range ue_info.Websites {
				if existing_website == new_website {
					return &protos.AddUEHeaderEnrichmentResult{Result: protos.AddUEHeaderEnrichmentResult_WEBSITE_CONFLICT}, nil
				}
			}
		}
	}

	s.ueInfos[ueIp][req.RuleId] = &control_plane.UEInfo{
		Websites: req.Websites,
		Headers:  req.Headers,
	}

	s.controllerCli.UpdateSnapshot(s.ueInfos)

	return &protos.AddUEHeaderEnrichmentResult{Result: protos.AddUEHeaderEnrichmentResult_SUCCESS}, nil
}

// DeactivateUEHeaderEnrichment deactivates/removes the UE from the current header enrichment list
func (s *envoyControllerService) DeactivateUEHeaderEnrichment(
	ctx context.Context,
	req *protos.DeactivateUEHeaderEnrichmentRequest,
) (*protos.DeactivateUEHeaderEnrichmentResult, error) {
	glog.Infof("DeactivateUEHeaderEnrichmentResult received for IP %s", req.UeIp.Address)
	glog.V(2).Infof("req %s", (req))

	ueIp := string(req.UeIp.Address)
	if _, ok := s.ueInfos[ueIp]; !ok {
		return &protos.DeactivateUEHeaderEnrichmentResult{Result: protos.DeactivateUEHeaderEnrichmentResult_UE_NOT_FOUND}, nil
	}
	if req.RuleId != "" {
		if _, ok := s.ueInfos[ueIp][req.RuleId]; !ok {
			return &protos.DeactivateUEHeaderEnrichmentResult{Result: protos.DeactivateUEHeaderEnrichmentResult_RULE_NOT_FOUND}, nil
		}
		delete(s.ueInfos[ueIp], req.RuleId)
		if len(s.ueInfos[ueIp]) == 0 {
			delete(s.ueInfos, ueIp)
		}
	} else {
		delete(s.ueInfos, ueIp)
	}
	s.controllerCli.UpdateSnapshot(s.ueInfos)

	return &protos.DeactivateUEHeaderEnrichmentResult{Result: protos.DeactivateUEHeaderEnrichmentResult_SUCCESS}, nil
}

// NewenvoyControllerService returns a new EnvoyController service
func NewEnvoyControllerService(controllerCli control_plane.EnvoyController) protos.EnvoyControllerServer {
	return &envoyControllerService{ueInfos: control_plane.UEInfoMap{}, controllerCli: controllerCli}
}
