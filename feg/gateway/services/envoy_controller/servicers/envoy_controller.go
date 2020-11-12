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
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/envoy_controller/control_plane"
	lte_proto "magma/lte/cloud/go/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

type envoyControllerService struct {
	ue_infos       []*protos.AddUEHeaderEnrichmentRequest
	controller_cli control_plane.EnvoyController
}

// AddUEHeaderEnrichment adds the UE to the current header enrichment list, if UE is already in the list replaces the he information for that UE
func (s *envoyControllerService) AddUEHeaderEnrichment(
	ctx context.Context,
	req *protos.AddUEHeaderEnrichmentRequest,
) (*protos.AddUEHeaderEnrichmentResult, error) {
	var (
		res *protos.AddUEHeaderEnrichmentResult
		err error
	)
	glog.Infof("AddUEHeaderEnrichmentResult received for IP %s", req.UeIp.Address)
	s.ue_infos = add_to_list(s.ue_infos, req)
	s.controller_cli.UpdateSnapshot(s.ue_infos)

	return res, err
}

// DeactivateUEHeaderEnrichment deactivates/removes the UE from the current header enrichment list
func (s *envoyControllerService) DeactivateUEHeaderEnrichment(
	ctx context.Context,
	req *protos.DeactivateUEHeaderEnrichmentRequest,
) (*protos.DeactivateUEHeaderEnrichmentResult, error) {
	var (
		res *protos.DeactivateUEHeaderEnrichmentResult
		err error
	)
	glog.Infof("DeactivateUEHeaderEnrichmentResult received for IP %s", req.UeIp.Address)
	s.ue_infos = remove_from_list(s.ue_infos, req.UeIp)
	s.controller_cli.UpdateSnapshot(s.ue_infos)

	return res, err
}

// NewenvoyControllerService returns a new EnvoyController service
func NewEnvoyControllerService(controller_cli control_plane.EnvoyController) protos.EnvoyControllerServer {
	return &envoyControllerService{controller_cli: controller_cli}
}

func remove_from_list(l []*protos.AddUEHeaderEnrichmentRequest, ip *lte_proto.IPAddress) []*protos.AddUEHeaderEnrichmentRequest {
	for i, other := range l {
		if string(other.UeIp.Address) == string(ip.Address) {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func add_to_list(l []*protos.AddUEHeaderEnrichmentRequest, new *protos.AddUEHeaderEnrichmentRequest) []*protos.AddUEHeaderEnrichmentRequest {
	for i, other := range l {
		if string(other.UeIp.Address) == string(new.UeIp.Address) {
		    // Overwrite duplicate UE
		    ret := append(l[:i], l[i+1:]...)
		    return append(ret, new)
		}
	}
	l = append(l, new)
	return l
}
