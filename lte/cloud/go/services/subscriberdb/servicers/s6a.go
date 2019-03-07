package servicers

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/protos"
)

func (srv *SubscriberDBServer) AuthenticationInformation(ctx context.Context, air *protos.AuthenticationInformationRequest) (*protos.AuthenticationInformationAnswer, error) {
	return nil, status.Errorf(codes.Unimplemented, "authentication information not implemented")
}

func (srv *SubscriberDBServer) UpdateLocation(ctx context.Context, ulr *protos.UpdateLocationRequest) (*protos.UpdateLocationAnswer, error) {
	return nil, status.Errorf(codes.Unimplemented, "update location not implemented")
}

func (srv *SubscriberDBServer) PurgeUE(ctx context.Context, purge *protos.PurgeUERequest) (*protos.PurgeUEAnswer, error) {
	return nil, status.Errorf(codes.Unimplemented, "purge UE not implemented")
}
