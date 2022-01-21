package subscriberdb

import (
	"context"

	"github.com/magma/magma/src/go/log"
	"github.com/magma/magma/src/go/protos/magma/orc8r"
	pb "github.com/magma/magma/src/go/protos/magma/subscriberdb"
	"github.com/magma/magma/src/go/subscriberdb"
)

type SubscriberDBServer struct {
	log.Logger
	pb.SubscriberDBServer
	subscriberdb.SubscriberDB
}

func NewSubscriberDBServer(logger log.Logger, db subscriberdb.SubscriberDB) *SubscriberDBServer {
	return &SubscriberDBServer{
		Logger:       logger,
		SubscriberDB: db,
	}
}

func (s *SubscriberDBServer) AddSubscriber(ctx context.Context, req *pb.SubscriberData) (*orc8r.Void, error) {
	s.Logger.With("subscriberdb", req).Debug().Print("AddSubscriber")
	err := s.SubscriberDB.Add(req.GetSid().GetId(), *req)
	if err != nil {
		return nil, err
	}
	return &orc8r.Void{}, nil
}

func (s *SubscriberDBServer) DeleteSubscriber(ctx context.Context, req *pb.SubscriberID) (*orc8r.Void, error) {
	s.Logger.With("subscriberdb", req).Debug().Print("DeleteSubscriber")
	err := s.SubscriberDB.Delete(req.GetId())
	if err != nil {
		return nil, err
	}
	return &orc8r.Void{}, nil
}

func (s *SubscriberDBServer) UpdateSubscriber(ctx context.Context, req *pb.SubscriberUpdate) (*orc8r.Void, error) {
	s.Logger.With("subscriberdb", req).Debug().Print("UpdateSubscriber")
	err := s.SubscriberDB.Update(req.GetData().GetSid().GetId(), *req.Data)
	if err != nil {
		return nil, err
	}
	return &orc8r.Void{}, nil
}

func (s *SubscriberDBServer) GetSubscriberData(ctx context.Context, req *pb.SubscriberID) (*pb.SubscriberData, error) {
	s.Logger.With("subscriberdb", req).Debug().Print("GetSubscriberData")
	data, err := s.SubscriberDB.Get(req.GetId())
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *SubscriberDBServer) ListSubscribers(ctx context.Context, req *orc8r.Void) (*pb.SubscriberIDSet, error) {
	s.Logger.With("subscriberdb", req).Debug().Print("ListSubscribers")
	subscribers, err := s.List()
	if err != nil {
		return nil, err
	}
	return &pb.SubscriberIDSet{Sids: subscribers}, nil
}
