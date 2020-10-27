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

// Package servicers implements Swx GRPC proxy service which sends MAR/SAR messages over
// diameter connection, waits (blocks) for diameter's MAA/SAAs returns their RPC representation
package servicers

import (
	"github.com/golang/glog"
	"golang.org/x/net/context"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/envoy_controller/control_plane"
	lte_proto "magma/lte/cloud/go/protos"
)

const (
	TIMEOUT = 10
)

type envoydService struct {
	ue_infos []*protos.AddUEHeaderEnrichmentRequest
}

// Register sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s *envoydService) AddUEHeaderEnrichment(
	ctx context.Context,
	req *protos.AddUEHeaderEnrichmentRequest,
) (*protos.AddUEHeaderEnrichmentResult, error) {
	var (
		res *protos.AddUEHeaderEnrichmentResult
		err error
	)
	s.ue_infos = append(s.ue_infos, req)

	glog.Infof("AddUEHeaderEnrichmentResult received")
	control_plane.UpdateSnapshot(s.ue_infos)

	return res, err
}

func (s *envoydService) DeactivateUEHeaderEnrichment(
	ctx context.Context,
	req *protos.DeactivateUEHeaderEnrichmentRequest,
) (*protos.DeactivateUEHeaderEnrichmentResult, error) {
	var (
		res *protos.DeactivateUEHeaderEnrichmentResult
		err error
	)
	glog.Infof("DeactivateUEHeaderEnrichmentResult received")
	s.ue_infos = remove(s.ue_infos, req.UeIp)
	return res, err
}

// NewenvoydService returns a new Envoyd service
func NewEnvoydService() protos.EnvoydServer {
	return &envoydService{}
}

func remove(l []*protos.AddUEHeaderEnrichmentRequest, ip *lte_proto.IPAddress) []*protos.AddUEHeaderEnrichmentRequest {
	for i, other := range l {
		if string(other.UeIp.Address) == string(ip.Address) {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}
