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

// package service implements magmad GRPC service
package service

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/aeden/traceroute"
	"github.com/emakeev/snowflake"
	"github.com/golang/glog"

	"magma/gateway/config"
	"magma/gateway/mconfig"
	config_service "magma/gateway/services/configurator/service"
	"magma/gateway/services/magmad/service/generic_command"
	"magma/gateway/services/magmad/service/ping"
	"magma/gateway/services/magmad/service_manager"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/service"
)

type magmadService struct {
	protos.UnimplementedMagmadServer
}

func (m *magmadService) StartServices(context.Context, *protos.Void) (*protos.Void, error) {
	resErrs := errors.NewMulti()
	sm := service_manager.Get()
	for _, srv := range getServices() {
		glog.Infof("Starting service '%s'", srv)
		resErrs = resErrs.AddFmt(sm.Start(srv), "service '%s' start error:", srv)
	}
	return &protos.Void{}, resErrs.AsError()
}

func (m *magmadService) StopServices(context.Context, *protos.Void) (*protos.Void, error) {
	resErrs := errors.NewMulti()
	sm := service_manager.Get()
	for _, srv := range getServices() {
		glog.Infof("Stopping service '%s'", srv)
		resErrs = resErrs.AddFmt(sm.Stop(srv), "service '%s' stop error:", srv)
	}
	return &protos.Void{}, resErrs.AsError()
}

func (m *magmadService) Reboot(context.Context, *protos.Void) (*protos.Void, error) {
	glog.Info("Rebooting Gateway")
	go exec.Command("reboot").Run()
	return &protos.Void{}, nil
}

func (m *magmadService) RestartServices(context.Context, *protos.RestartServicesRequest) (*protos.Void, error) {
	resErrs := errors.NewMulti()
	sm := service_manager.Get()
	for _, srv := range getServices() {
		glog.Infof("Restarting service '%s'", srv)
		resErrs = resErrs.AddFmt(sm.Restart(srv), "service '%s' restart error:", srv)
	}
	return &protos.Void{}, resErrs.AsError()
}

func (m *magmadService) GetConfigs(context.Context, *protos.Void) (*protos.GatewayConfigs, error) {
	return mconfig.GetGatewayConfigs(), nil
}

func (m *magmadService) SetConfigs(_ context.Context, cfg *protos.GatewayConfigs) (*protos.Void, error) {
	var err error
	if cfg != nil {
		var marshaled []byte
		marshaled, err = protos.MarshalMconfig(cfg)
		if err == nil {
			err = config_service.SaveConfigs(marshaled)
		}
	}
	return &protos.Void{}, err
}

func (m *magmadService) RunNetworkTests(ctx context.Context, req *protos.NetworkTestRequest) (*protos.NetworkTestResponse, error) {
	res := &protos.NetworkTestResponse{}
	if req == nil {
		return res, nil
	}
	execCtx, cancel := context.WithTimeout(ctx, generic_command.Timeout)
	defer cancel()
	// Process pings
	for _, png := range req.Pings {
		if png == nil {
			continue
		}
		pingRes := &protos.PingResult{HostOrIp: png.HostOrIp, NumPackets: 4}
		if png.NumPackets > 0 {
			pingRes.NumPackets = png.NumPackets
		}
		packets := strconv.FormatInt(int64(pingRes.NumPackets), 10)
		cmd := exec.CommandContext(execCtx, "ping", "-c", packets, png.HostOrIp)
		glog.Info(cmd.String())
		out, err := cmd.Output()
		if err != nil {
			pingRes.Error = fmt.Sprintf("error executing '%s': %v", cmd.String(), err)
		}
		if len(out) > 0 {
			ping.ParseResult(out, pingRes)
		}
		res.Pings = append(res.Pings, pingRes)
	}
	// Process traceroutes
	for _, trt := range req.Traceroutes {
		if trt == nil {
			continue
		}
		options := &traceroute.TracerouteOptions{}
		options.SetMaxHops(int(trt.GetMaxHops()))
		options.SetPacketSize(int(trt.GetBytesPerPacket()))
		glog.Infof("traceroute %s -m %d %d", trt.GetHostOrIp(), options.MaxHops(), options.PacketSize())
		trtRes := &protos.TracerouteResult{
			HostOrIp: trt.GetHostOrIp(),
		}
		hc := make(chan traceroute.TracerouteHop)
		// start 'streaming' the chan
		go func() {
			hopMap := map[int]*protos.TracerouteHop{}
			for hop := range hc {
				var probe *protos.TracerouteProbe
				if hop.Success {
					probe = &protos.TracerouteProbe{
						Hostname: hop.HostOrAddressString(),
						Ip:       hop.AddressString(),
						RttMs:    Milliseconds(hop.ElapsedTime),
					}
				} else {
					probe = &protos.TracerouteProbe{
						Hostname: "*",
						Ip:       "*",
					}
				}
				hopRes, ok := hopMap[hop.TTL]
				if ok && hopRes != nil {
					hopRes.Probes = append(hopRes.Probes, probe)
				} else {
					hopRes = &protos.TracerouteHop{
						Idx:    int32(hop.TTL),
						Probes: []*protos.TracerouteProbe{probe},
					}
					trtRes.Hops = append(trtRes.Hops, hopRes)
					hopMap[hop.TTL] = hopRes
				}
			}
		}()
		_, err := traceroute.Traceroute(trt.GetHostOrIp(), options, hc)
		// the last write & close of hc chan should complete before Traceroute() returns
		if err != nil {
			trtRes.Error = err.Error()
		}
		res.Traceroutes = append(res.Traceroutes, trtRes)
	}
	return res, nil
}

// Seconds returns the duration as a floating point number of seconds.
func Milliseconds(d time.Duration) float32 {
	ms := d / time.Millisecond
	nsec := d % time.Millisecond
	return float32(ms) + float32(nsec)/float32(time.Millisecond)
}

func (m *magmadService) GenericCommand(
	ctx context.Context, req *protos.GenericCommandParams) (*protos.GenericCommandResponse, error) {

	return generic_command.Execute(ctx, req)
}

func (m *magmadService) GetGatewayId(context.Context, *protos.Void) (*protos.GetGatewayIdResponse, error) {
	id, err := snowflake.Get()
	resp := &protos.GetGatewayIdResponse{GatewayId: id.String()}
	return resp, err
}

func (m *magmadService) TailLogs(req *protos.TailLogsRequest, srv protos.Magmad_TailLogsServer) error {
	c, proc, err := service_manager.Get().TailLogs(req.GetService())
	if err != nil {
		return err
	}
	for s := range c {
		err := srv.SendMsg(&protos.LogLine{Line: s + "\n"}) // append LF, our existing tools rely on it for printouts
		if err != nil {
			proc.Kill()
			return err
		}
	}
	return nil
}

// NewMagmadService returns a new magmad service
func NewMagmadService() protos.MagmadServer {
	return &magmadService{}
}

// StartMagmadServer runs instance of the magmad grpc service
// StartMagmadServer only returns on error and has to be run in its own Go routine or main thread
func StartMagmadServer() error {
	srv, err := service.NewServiceWithOptions("", strings.ToUpper(definitions.MagmadServiceName))
	if err != nil {
		return fmt.Errorf("error creating '%s' service: %v", definitions.MagmadServiceName, err)
	}
	protos.RegisterMagmadServer(srv.GrpcServer, NewMagmadService())
	glog.Infof("starting '%s' Service", definitions.MagmadServiceName)
	err = srv.Run()
	if err != nil {
		return fmt.Errorf("error starting '%s' service: %s", definitions.MagmadServiceName, err)
	}
	return nil
}

func getServices() []string {
	mdc := config.GetMagmadConfigs()
	return mdc.MagmaServices
}
