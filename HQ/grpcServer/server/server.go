package grpcServer

import (
	"context"
	"io"
	"log"
	"time"

	db "github.com/NaduniRabel/distributed-system-monitor/HQ/db"
	pb "github.com/NaduniRabel/distributed-system-monitor/proto"
	"github.com/jackc/pgx/v5/pgxpool"
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
		pool := db.Pool
		SaveData(req, pool)
	}
}

func SaveData(data *pb.Metrics, pool *pgxpool.Pool) error {
	ctx := context.Background()
	services := data.Services
	systemQuery := `INSERT INTO system_logs (
				host_id, 
				recorded_at, 
				cpu_usage, 
				total_ram, 
				free_ram, 
				used_ram, 
				available_ram, 
				total_disk, 
				used_disk, 
				free_disk
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
			) 
			ON CONFLICT (host_id, recorded_at) 
			DO UPDATE SET 
				cpu_usage     = EXCLUDED.cpu_usage,
				free_ram      = EXCLUDED.free_ram,
				used_ram      = EXCLUDED.used_ram,
				available_ram = EXCLUDED.available_ram,
				used_disk     = EXCLUDED.used_disk,
				free_disk     = EXCLUDED.free_disk;`

	_, err := pool.Exec(ctx, systemQuery,
		data.Host.HostID,
		time.Now(),
		data.CPU.Usage[0],
		data.Memory.Total,
		data.Memory.Free,
		data.Memory.UsedPercentage,
		data.Memory.Available,
		data.Disk.Total,
		data.Disk.Used,
		data.Disk.Free,
	)

	serviceQuery := `
		INSERT INTO service_logs (
			host_id, 
			service_name, 
			recorded_at, 
			status, 
			cpu_usage, 
			ram_usage
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) 
		ON CONFLICT (host_id, recorded_at, service_name) 
		DO UPDATE SET 
			status    = EXCLUDED.status,
			cpu_usage = EXCLUDED.cpu_usage,
			ram_usage = EXCLUDED.ram_usage;`

	if err != nil {
		log.Printf("Error saving system log %e", err)
	} else {
		log.Print("System data successfully saved")
	}

	for _, service := range services {

		_, err := pool.Exec(ctx, serviceQuery,
			data.Host.HostID,
			service.Name,
			time.Now(),
			service.Status,
			service.CPU,
			service.Memory,
		)

		if err != nil {
			log.Printf("An error occured while saving %v", service.Name)
			continue
		}
	}

	return nil
}
