package main

import (
	"context"
	"errors"
	"log"

	pb "github.com/aryan9600/delivery-system/service-vessel/proto/vessel"
	"github.com/micro/go-micro/v2"
)

type Repository interface {
	FindVessel(*pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
	vessels []*pb.Vessel
}

func (repo *VesselRepository) FindVessel(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity < vessel.Capacity && spec.MaxWeight < vessel.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("No vessel found for this spec")
}

type vesselService struct {
	repo Repository
}

func (s *vesselService) FindFree(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	vessel, err := s.repo.FindVessel(req)
	if err != nil {
		return err
	}
	res.Vessel = vessel
	return nil
}

func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{
			Id:        "v001",
			Name:      "USSR",
			MaxWeight: 6900,
			Capacity:  20,
		},
	}
	repo := &VesselRepository{vessels: vessels}
	service := micro.NewService(micro.Name("service.vessel"))
	service.Init()

	if err := pb.RegisterVesselServiceHandler(service.Server(), &vesselService{repo}); err != nil {
		log.Panic(err)
	}
	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
