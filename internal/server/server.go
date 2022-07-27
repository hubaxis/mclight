package server

import (
	"context"
	"github.com/hubaxis/mclight/internal/cache"
	"github.com/hubaxis/mclight/protocol/mclight"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//Server encapsulates GRPC protocol
type Server struct {
	mclight.UnimplementedMCLightServiceServer
	cache cache.Cache
}

// New memcached instance
func New(cache cache.Cache) *Server {
	return &Server{cache: cache}
}

// Set data to memcached
func (s Server) Set(_ context.Context, request *mclight.SetRequest) (*mclight.SetResponse, error) {
	err := s.cache.Set(request.Key, request.Value, request.Expiration.AsDuration())
	if err != nil {
		log.Errorf("can't set value %v %v", request, err)
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &mclight.SetResponse{}, nil
}

// Delete data from memcached
func (s Server) Delete(_ context.Context, request *mclight.DeleteRequest) (*mclight.DeleteResponse, error) {
	err := s.cache.Delete(request.Key)
	if err != nil {
		log.Errorf("can't delete value %v %v", request, err)
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &mclight.DeleteResponse{}, nil
}

// Get data from memcached
func (s Server) Get(_ context.Context, request *mclight.GetRequest) (*mclight.GetResponse, error) {
	value, err := s.cache.Get(request.Key)
	if err != nil {
		log.Errorf("can't get value %v %v", request, err)
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &mclight.GetResponse{Value: value}, nil
}
