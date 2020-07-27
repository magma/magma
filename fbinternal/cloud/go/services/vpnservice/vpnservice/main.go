package main

import (
	"log"

	"magma/fbinternal/cloud/go/fbinternal"
	"magma/fbinternal/cloud/go/protos"
	"magma/fbinternal/cloud/go/services/vpnservice"
	"magma/fbinternal/cloud/go/services/vpnservice/servicers"
	"magma/orc8r/cloud/go/service"
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
