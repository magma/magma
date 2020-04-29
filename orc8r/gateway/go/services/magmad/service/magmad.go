/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package service implements magmad GRPC service
package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/aeden/traceroute"
	"github.com/emakeev/snowflake"
	
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
}

func (m *magmadService) StartServices(context.Context, *protos.Void) (*protos.Void, error) {
	var resErrs *errors.MultiError
	sm := service_manager.Get()
	for _, srv := range getServices() {
		resErrs.Add(sm.Start(srv))
	}
	return &protos.Void{}, resErrs
}

func (m *magmadService) StopServices(context.Context, *protos.Void) (*protos.Void, error) {
	var resErrs *errors.MultiError
	sm := service_manager.Get()
	for _, srv := range getServices() {
		resErrs.Add(sm.Stop(srv))
	}
	return &protos.Void{}, resErrs
}

func (m *magmadService) Reboot(context.Context, *protos.Void) (*protos.Void, error) {
	go exec.Command("reboot").Run()
	return &protos.Void{}, nil
}

func (m *magmadService) RestartServices(context.Context, *protos.RestartServicesRequest) (*protos.Void, error) {
	var resErrs *errors.MultiError
	sm := service_manager.Get()
	for _, srv := range getServices() {
		resErrs.Add(sm.Restart(srv))
	}
	return &protos.Void{}, resErrs
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
			_, err = config_service.SaveConfigs(marshaled, false)
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
		trtRes := &protos.TracerouteResult{
			HostOrIp: trt.GetHostOrIp(),
		}
		maxHops := options.MaxHops()
		hc := make(chan traceroute.TracerouteHop)
		// start 'streaming' the chan
		go func() {
			for hop := range hc {
				var ipStr, hostStr string
				if hop.Success {
					ipStr = net.IP(hop.Address[:]).String()
					hostStr = hop.Host
				}
				trtRes.Hops = append(trtRes.Hops, &protos.TracerouteHop{
					Idx: int32(maxHops - hop.TTL),
					Probes: []*protos.TracerouteProbe{{
						Hostname: hostStr,
						Ip:       ipStr,
						RttMs:    Milliseconds(hop.ElapsedTime),
					}},
				})
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
		err := srv.SendMsg(&protos.LogLine{Line: s})
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
	log.Printf("starting '%s' Service", definitions.MagmadServiceName)
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
