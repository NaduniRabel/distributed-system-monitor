package grpcServer

import (
	"io"
	"log"

	pb "github.com/NaduniRabel/distributed-system-monitor/proto"
)

/*Struct*/
type MetricServiceServer struct {
	pb.UnimplementedMetricServiceServer
}

/*Receive the stream*/
func (s *MetricServiceServer) StreamMetrics(stream pb.MetricService_StreamMetricsServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Connection closed by client")
			return nil
		}
		if err != nil {
			return err
		}

		metrics := req.Host.HostID

		for _, metric := range metrics {
			log.Println("Message received from ", metric)
		}

	}
}
