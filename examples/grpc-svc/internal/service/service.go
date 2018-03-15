package service

import (
	"context"

	"github.com/utilitywarehouse/go-circle-build/examples/grpc-svc/pkg/pb/example"
)

// SVC implements the example interface and allows embedding additional code
type SVC struct {
	db string // fake db
}

// New returns a service with embedded etcd client plus
func New(db string) *SVC {
	return &SVC{db}
}

// GetExample does nothing
func (s *SVC) GetExample(ctx context.Context, req *example.ExampleRequest) (*example.ExampleResponse, error) {
	resp := req.GetReq() + s.db
	return &example.ExampleResponse{Resp: resp}, nil
}
