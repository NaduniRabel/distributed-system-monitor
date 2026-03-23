package main

import (
	"log"
	"net"

	db "github.com/NaduniRabel/distributed-system-monitor/HQ/db"
	grpcServer "github.com/NaduniRabel/distributed-system-monitor/HQ/grpcServer/server"
	pb "github.com/NaduniRabel/distributed-system-monitor/proto"
	"google.golang.org/grpc"
)

func main() {

	db.Init()

	/*Listen on port 9001*/
	lis, err := net.Listen("tcp", ":9001")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	} else {
		log.Println("Listening on port 9001")
	}
	/*Create new server*/
	grpcServ := grpc.NewServer()

	log.Println("Registering server...")

	/*Register server*/
	pb.RegisterMetricServiceServer(grpcServ, &grpcServer.MetricServiceServer{})

	log.Println("Server registered successfully")

	err = grpcServ.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve %s", err)
	}

}
