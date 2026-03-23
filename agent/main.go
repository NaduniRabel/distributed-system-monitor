package main

import (
	"context"
	"log"
	"time"

	"github.com/NaduniRabel/distributed-system-monitor/agent/collector"
	pb "github.com/NaduniRabel/distributed-system-monitor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	/*Connect with the gRPC server running on 9001*/
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMetricServiceClient(conn)

	/*Initialize the stream*/
	stream, err := client.StreamMetrics(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	log.Println("Agent started. Streaming metrics every 5 seconds...")

	/*Ticker with 5-second intervals*/
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		/*Collect emtrics*/
		data, err := collector.CollectMetrics()
		if err != nil {
			log.Printf("Failed to collect metrics: %v", err)
			continue
		}

		/*Map collector struct to Protobuf message*/
		protoMetrics := &pb.Metrics{
			Host: &pb.HostMetrics{HostID: data.Host.HostID},
			CPU:  &pb.CPUMetrics{Usage: data.CPU.Usage},
			Memory: &pb.MemMetrics{
				Total:          data.Memory.Total,
				Free:           data.Memory.Free,
				UsedPercentage: data.Memory.UsedPercentage,
				Available:      data.Memory.Available,
			},
			Disk: &pb.DiskMetrics{
				Total: data.Disk.Total,
				Used:  data.Disk.Used,
				Free:  data.Disk.Free,
			},
		}

		for _, s := range data.Services {
			protoMetrics.Services = append(protoMetrics.Services, &pb.ServiceMetrics{
				Name:   s.Name,
				Status: s.Status,
				CPU:    s.CPU,
				Memory: s.Memory,
			})
		}

		/* Send data*/
		if err := stream.Send(protoMetrics); err != nil {
			log.Printf("Failed to send metrics over stream: %v", err)
		} else {
			log.Printf("Successfully streamed metrics for Host: %s", data.Host.HostID)
		}
	}

}
