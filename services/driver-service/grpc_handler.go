package main

import (
	"context"
	"fmt"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedDriverServiceServer

	service *Service
}

func NewGRPCHandler(server *grpc.Server, service *Service) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}

	pb.RegisterDriverServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	// TODO: Call the service method
	driverID := req.GetDriverID()
	packageSlug := req.GetPackageSlug()

	driver, err := h.service.RegisterDriver(driverID, packageSlug)
	if err != nil {
		return nil, fmt.Errorf("error registering driver: %v", err)
	}

	return &pb.RegisterDriverResponse{
		Driver: driver,
	}, nil
}

func (h *gRPCHandler) UnregisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	// TODO: Call the service method
	driverID := req.GetDriverID()
	h.service.UnregisterDriver(driverID)
	return &pb.RegisterDriverResponse{
			Driver: &pb.Driver{
				Id: driverID,
			},
		},
		nil
}
