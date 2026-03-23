package main

import (
	"github.com/NaduniRabel/distributed-system-monitor/agent/collector"
	//pb "github.com/NaduniRabel/distributed-system-monitor/proto"
)

func main() {

	collector.CollectMetrics()
	//pb.NewMetricServiceClient()

}
