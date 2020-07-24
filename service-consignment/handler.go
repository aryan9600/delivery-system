package main

import (
	"context"
	"errors"

	pb "github.com/aryan9600/delivery-system/service-consignment/proto/consignment"
	vessel "github.com/aryan9600/delivery-system/service-vessel/proto/vessel"
)

type handler struct {
	repository
	vesselClient vessel.VesselService
}

func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	vesselResponse, err := s.vesselClient.FindFree(ctx, &vessel.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if vesselResponse == nil {
		return errors.New("Something went wrong while fetching vessel.")
	}
	if err != nil {
		return err
	}
	req.VesselId = vesselResponse.Vessel.Id

	if err := s.repository.Create(ctx, MarshalConsignment(req)); err != nil {
		return err
	}
	res.Created = true
	res.Consignment = req
	return nil
}

func (s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments, err := s.repository.GetAll(ctx)
	if err != nil {
		return err
	}
	res.Consignments = UnmarshalConsignemntCollection(consignments)
	return nil
}
