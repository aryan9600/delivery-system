package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/aryan9600/delivery-system/service-consignment/proto/consignment"
	vessel "github.com/aryan9600/delivery-system/service-vessel/proto/vessel"
	"github.com/micro/go-micro/v2"
)

const (
	defaultHost = "datastore:27017"
)

func main() {
	service := micro.NewService(micro.Name("service.cosnignment"))
	service.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}
	client, err := CreateClient(context.Background(), uri, 2)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())

	consginmentCollection := client.Database("delivery").Collection("consignments")
	repository := &MongoRepistory{consginmentCollection}
	vesselClient := vessel.NewVesselService("service.client", service.Client())

	h := &handler{repository, vesselClient}
	pb.RegisterShippingServiceHandler(service.Server(), h)
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
