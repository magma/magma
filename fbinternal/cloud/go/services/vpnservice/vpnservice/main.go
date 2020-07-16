package main

import (
	"log"

	"magma/orc8r/cloud/go/service"
	"orc8r/fbinternal/cloud/go/fbinternal"
	"orc8r/fbinternal/cloud/go/protos"
	"orc8r/fbinternal/cloud/go/services/vpnservice"
	"orc8r/fbinternal/cloud/go/services/vpnservice/servicers"
)

const taKeyPath = "/var/opt/magma/certs/vpn_ta.key"

func main() {
	srv, err := service.NewOrchestratorService(fbinternal.ModuleName, vpnservice.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	// Add servicers to the service
	servicer := servicers.NewVPNServicer(taKeyPath)
	protos.RegisterVPNServiceServer(srv.GrpcServer, servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
