package main

import (
	"context"

	pb "github.com/aryan9600/delivery-system/service-consignment/proto/consignment"
	"go.mongodb.org/mongo-driver/mongo"
)

type Consignment struct {
	ID          string       `json:"id"`
	Weight      int32        `json:"weight"`
	Description string       `json:"description"`
	Containers  []*Container `json:"containers"`
	VesselID    string       `json:"vessel_id"`
}

type Container struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	UserID     string `json:"user_id"`
}

func MarshalContainerCollection(containers []*pb.Container) []*Container {
	collection := make([]*Container, 0)
	for _, container := range containers {
		collection = append(collection, MarshalContainer(container))
	}
	return collection
}

func UnmarshalContainerCollection(containers []*Container) []*pb.Container {
	collection := make([]*pb.Container, 0)
	for _, container := range containers {
		collection = append(collection, UnmarshalContainer(container))
	}
	return collection
}

func UnmarshalConsignemntCollection(consignments []*Consignment) []*pb.Consignment {
	collection := make([]*pb.Consignment, 0)
	for _, consignment := range consignments {
		collection = append(collection, UnmarshalConsignment(consignment))
	}
	return collection
}

func MarshalContainer(container *pb.Container) *Container {
	return &Container{
		ID:         container.Id,
		CustomerID: container.CustomerId,
		UserID:     container.UserId,
	}
}

func UnmarshalContainer(container *Container) *pb.Container {
	return &pb.Container{
		Id:         container.ID,
		UserId:     container.UserID,
		CustomerId: container.CustomerID,
	}
}

func UnmarshalConsignment(consignment *Consignment) *pb.Consignment {
	return &pb.Consignment{
		Id:          consignment.ID,
		Description: consignment.Description,
		Weight:      consignment.Weight,
		Containers:  UnmarshalContainerCollection(consignment.Containers),
		VesselId:    consignment.VesselID,
	}
}

func MarshalConsignment(consignment *pb.Consignment) *Consignment {
	return &Consignment{
		ID:          consignment.Id,
		Weight:      consignment.Weight,
		Description: consignment.Description,
		Containers:  MarshalContainerCollection(consignment.Containers),
		VesselID:    consignment.VesselId,
	}
}

type repository interface {
	Create(ctx context.Context, consignment *Consignment) error
	GetAll(ctx context.Context) ([]*Consignment, error)
}

type MongoRepistory struct {
	collection *mongo.Collection
}

func (repo *MongoRepistory) Create(ctx context.Context, consignment *Consignment) error {
	_, err := repo.collection.InsertOne(ctx, consignment)
	return err
}

func (repo *MongoRepistory) GetAll(ctx context.Context) ([]*Consignment, error) {
	cur, err := repo.collection.Find(ctx, nil, nil)
	var consignments []*Consignment
	for cur.Next(ctx) {
		var consignment *Consignment
		if err := cur.Decode(&consignment); err != nil {
			return nil, err
		}
		consignments = append(consignments, consignment)
	}
	return consignments, err
}
