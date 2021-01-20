/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers

import (
	"errors"
	"net/http"
	"sort"

	"magma/fbinternal/cloud/go/services/testcontroller"
	"magma/fbinternal/cloud/go/services/testcontroller/obsidian/models"
	"magma/fbinternal/cloud/go/services/testcontroller/storage"
	"magma/orc8r/cloud/go/obsidian"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/labstack/echo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CIRoot        = "ci"
	CINodesBase   = "nodes"
	CIReserveBase = "reserve"
	CIReleaseBase = "release"

	NodeIDArg  = ":node_id"
	LeaseIDArg = ":lease_id"

	CIRootPath                 = obsidian.V1Root + CIRoot
	CINodesRootPath            = CIRootPath + obsidian.UrlSep + CINodesBase
	CINodesGetPath             = CINodesRootPath + obsidian.UrlSep + NodeIDArg
	CINodesReservePath         = CIRootPath + obsidian.UrlSep + CIReserveBase
	CINodesManuallyReservePath = CINodesGetPath + obsidian.UrlSep + CIReserveBase
	CINodesManuallyReleasePath = CINodesGetPath + obsidian.UrlSep + CIReleaseBase
	CINodesReleasePath         = CINodesGetPath + obsidian.UrlSep + CIReleaseBase + obsidian.UrlSep + LeaseIDArg

	tagParamName          = "tag"
	listUntaggedParamName = "list_untagged"
)

func listCINodes(c echo.Context) error {
	tagVal := c.QueryParam(tagParamName)
	shouldListUntagged := len(c.QueryParam(listUntaggedParamName)) > 0
	var tag *string
	if tagVal != "" {
		tag = &tagVal
	}
	if shouldListUntagged {
		tag = strPtr("")
	}

	nodes, err := testcontroller.GetNodes(nil, tag)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := make([]*models.CiNode, 0, len(nodes))
	for _, node := range nodes {
		ret = append(ret, protoNodeToModel(node))
	}
	sort.Slice(ret, func(i, j int) bool { return *ret[i].ID < *ret[j].ID })
	return c.JSON(http.StatusOK, ret)
}

func getCINode(c echo.Context) error {
	idParam, nerr := obsidian.GetParamValues(c, "node_id")
	if nerr != nil {
		return nerr
	}

	nodes, err := testcontroller.GetNodes(idParam, nil)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	node, found := nodes[idParam[0]]
	if !found {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, protoNodeToModel(node))
}

func createCINode(c echo.Context) error {
	node := &models.MutableCiNode{}
	if err := c.Bind(node); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := node.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	err := testcontroller.CreateOrUpdateNode(&storage.MutableCINode{Id: *node.ID, Tag: node.Tag, VpnIP: string(*node.VpnIP)})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func updateCINode(c echo.Context) error {
	idParam, nerr := obsidian.GetParamValues(c, "node_id")
	if nerr != nil {
		return nerr
	}

	node := &models.MutableCiNode{}
	if err := c.Bind(node); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := node.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	if *node.ID != idParam[0] {
		return obsidian.HttpError(errors.New("payload ID does not match path param"), http.StatusBadRequest)
	}
	err := testcontroller.CreateOrUpdateNode(&storage.MutableCINode{Id: *node.ID, Tag: node.Tag, VpnIP: string(*node.VpnIP)})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteCINode(c echo.Context) error {
	idParam, nerr := obsidian.GetParamValues(c, "node_id")
	if nerr != nil {
		return nerr
	}

	err := testcontroller.DeleteNode(idParam[0])
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func leaseCINode(c echo.Context) error {
	lease, err := testcontroller.LeaseNode(c.QueryParam(tagParamName))
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if lease == nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, &models.NodeLease{
		ID:      swag.String(lease.Id),
		LeaseID: swag.String(lease.LeaseID),
		VpnIP:   ipv4Ptr(lease.VpnIP),
	})
}

func reserveCINode(c echo.Context) error {
	idParam, nerr := obsidian.GetParamValues(c, "node_id")
	if nerr != nil {
		return nerr
	}

	lease, err := testcontroller.ReserveNode(idParam[0])
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if lease == nil {
		return obsidian.HttpError(errors.New("Either the node is not known or it has already been reserved."), http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, &models.NodeLease{
		ID:      swag.String(lease.Id),
		LeaseID: swag.String(lease.LeaseID),
		VpnIP:   ipv4Ptr(lease.VpnIP),
	})
}

func returnManuallyReservedCINode(c echo.Context) error {
	idParam, nerr := obsidian.GetParamValues(c, "node_id")
	if nerr != nil {
		return nerr
	}

	// TODO: maybe expose this constant from the storage package?
	err := testcontroller.ReleaseNode(idParam[0], "manual")
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func releaseCINode(c echo.Context) error {
	params, nerr := obsidian.GetParamValues(c, "node_id", "lease_id")
	if nerr != nil {
		return nerr
	}
	nodeID, leaseID := params[0], params[1]
	err := testcontroller.ReleaseNode(nodeID, leaseID)
	if err == nil {
		return c.NoContent(http.StatusNoContent)
	}

	// Figure out if the error was due to bad params
	rpcErr, isRpcErr := status.FromError(err)
	if !isRpcErr {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	} else {
		switch rpcErr.Code() {
		case codes.InvalidArgument:
			return obsidian.HttpError(rpcErr.Err(), http.StatusBadRequest)
		default:
			return obsidian.HttpError(rpcErr.Err(), http.StatusInternalServerError)
		}
	}
}

func protoNodeToModel(n *storage.CINode) *models.CiNode {
	lastLeaseTime, err := ptypes.Timestamp(n.LastLeaseTime)
	if err != nil {
		// Don't make bad timestamp a return-blocker
		glog.Errorf("timestamp failed validation: %s", err)
	}
	return &models.CiNode{
		Available:     swag.Bool(n.Available),
		ID:            swag.String(n.Id),
		Tag:           n.Tag,
		LastLeaseTime: strfmt.DateTime(lastLeaseTime),
		VpnIP:         ipv4Ptr(n.VpnIp),
	}
}

func ipv4Ptr(s string) *strfmt.IPv4 {
	ip := strfmt.IPv4(s)
	return &ip
}

func strPtr(s string) *string {
	return &s
}
