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

package system_health

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/golang/glog"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/vishvananda/netlink"
)

// SystemHealth defines an interface to fetch system health and enable/disable
// functionality necessary for promotion/demotion from failovers
type SystemHealth interface {
	// GetSystemStats provides system level health stats
	GetSystemStats() (*SystemStats, error)

	// Disable allows the disabling of system level functionality. It is
	// up to implementors to determine specific functionality.
	Disable() error

	// Enable allows the enabling of system level functionality. It is
	// up to implementors to determine specific functionality.
	Enable() error
}

// SystemsStats define the metrics this provider will collect.
type SystemStats struct {
	CpuUtilPct float32
	MemUtilPct float32
}

// CWAGSystemHealthProvider defines a system health provider.
type CWAGSystemHealthProvider struct {
	trafficInterface string
	virtualIP        string
}

// NewCWAGSystemHealthProvider creates a new CWAGSystemHealthProvider.
func NewCWAGSystemHealthProvider(eth string, virtualIP string) (*CWAGSystemHealthProvider, error) {
	return &CWAGSystemHealthProvider{
		trafficInterface: eth,
		virtualIP:        virtualIP,
	}, nil
}

// GetSystemStats collects and return the stats defined in SystemStats.
func (c *CWAGSystemHealthProvider) GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{}
	cpuUtilPctArray, cpuErr := cpu.Percent(0, false)
	if cpuErr == nil && len(cpuUtilPctArray) == 1 {
		stats.CpuUtilPct = float32(cpuUtilPctArray[0]) / 100
	}
	virtualMem, vmErr := mem.VirtualMemory()
	if vmErr == nil {
		stats.MemUtilPct = float32(virtualMem.UsedPercent) / 100
	}
	if cpuErr != nil || vmErr != nil {
		return stats, fmt.Errorf("Error collecting system stats; CPU Result: %v, MEM Result: %v,", cpuErr, vmErr)
	}
	return stats, nil
}

// Enable adds the VIP to the configured interface and updates every route to
// use the VIP as the src IP.
func (c *CWAGSystemHealthProvider) Enable() error {
	iface, err := netlink.LinkByName(c.trafficInterface)
	if err != nil {
		return err
	}
	vipAddr, err := netlink.ParseAddr(c.virtualIP)
	if err != nil {
		return err
	}
	exists, err := c.doesAddrExistForInterface(iface, vipAddr)
	if err != nil {
		return err
	}
	if !exists {
		glog.V(1).Infof("Adding VIP address '%s' to interface %s", vipAddr.IP.String(), c.trafficInterface)
		err = netlink.AddrAdd(iface, vipAddr)
		if err != nil {
			return err
		}
	}
	// Ensure gRPC request doesn't timeout due to gratuitous arp
	go c.sendGratuitousArpReply(vipAddr.IP.String())

	return c.updateInterfaceRouteSrcIP(iface, *vipAddr)
}

// Disable removes the VIP from the configured interface and updates every
// route for to use the physical IP as the src IP.
func (c *CWAGSystemHealthProvider) Disable() error {
	iface, err := netlink.LinkByName(c.trafficInterface)
	if err != nil {
		return err
	}
	vipAddr, err := netlink.ParseAddr(c.virtualIP)
	if err != nil {
		return err
	}
	exists, err := c.doesAddrExistForInterface(iface, vipAddr)
	if err != nil {
		return err
	}
	if exists {
		glog.V(1).Infof("Removing VIP address '%s' from interface %s", vipAddr.IP.String(), c.trafficInterface)
		err = netlink.AddrDel(iface, vipAddr)
		if err != nil {
			return err
		}
	}
	addrs, err := netlink.AddrList(iface, netlink.FAMILY_V4)
	if err != nil {
		return err
	}
	if len(addrs) == 0 {
		return fmt.Errorf("No physical IP address found for interface")
	}
	physicalIp := addrs[0]
	return c.updateInterfaceRouteSrcIP(iface, physicalIp)
}

func (c *CWAGSystemHealthProvider) updateInterfaceRouteSrcIP(iface netlink.Link, newSrcIP netlink.Addr) error {
	routes, err := netlink.RouteList(iface, netlink.FAMILY_V4)
	if err != nil {
		return err
	}
	var routeErrors []string
	for _, route := range routes {
		if route.Src.String() == newSrcIP.IP.String() {
			continue
		}
		glog.V(1).Infof("Updating route (Dst: %s, Gw: %s) to use IP '%s' as src", route.Dst, route.Gw, newSrcIP.IP)
		route.Src = newSrcIP.IP
		err = netlink.RouteReplace(&route)
		if err != nil {
			routeErrors = append(routeErrors, err.Error())
		}
	}
	if len(routeErrors) > 0 {
		return fmt.Errorf("Encountered the following errors while updating routes:\n%s\n", strings.Join(routeErrors, "\n"))
	}
	return nil
}

func (c *CWAGSystemHealthProvider) sendGratuitousArpReply(vip string) error {
	arpCmd := exec.Command("arping", "-A", "-I", c.trafficInterface, vip, "-w", "3")
	output, err := arpCmd.CombinedOutput()
	if err != nil {
		glog.Errorf("Received error when sending gratuitous ARP reply: %s", string(output))
	}
	return err
}

func (c *CWAGSystemHealthProvider) doesAddrExistForInterface(iface netlink.Link, targetAddr *netlink.Addr) (bool, error) {
	addrs, err := netlink.AddrList(iface, netlink.FAMILY_V4)
	if err != nil {
		return false, err
	}
	for _, addr := range addrs {
		if addr.IP.String() == targetAddr.IP.String() {
			return true, nil
		}
	}
	return false, nil
}
