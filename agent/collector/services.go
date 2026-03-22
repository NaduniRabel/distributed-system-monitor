package collector

import (
	"encoding/json"

	"os"

	"github.com/shirou/gopsutil/v4/process"
)

/*Structs*/
type ServiceList struct {
	Services []string `json:"services"`
}

type ServiceMetrics struct {
	Name   string
	Status string
	CPU    float64
	Memory float32
}

/*Load and decode JSON data into struct*/
func LoadJSONData(path string) (ServiceList, error) {

	serviceList := ServiceList{}
	/*Load JSON*/
	file, err := os.ReadFile(path)
	if err != nil {
		return ServiceList{}, err
	}
	/*decode*/
	err = json.Unmarshal(file, &serviceList)
	if err != nil {
		return ServiceList{}, err
	}
	return serviceList, nil
}

/*Service level metrics*/
func GetServiceMetrics() ([]ServiceMetrics, error) {

	/*Get services from the JSON as a Struct*/
	services, err := LoadJSONData("config/services.json")
	if err != nil {
		return []ServiceMetrics{}, err
	}
	/*Get all processes*/
	processes, err := process.Processes()
	if err != nil {
		return []ServiceMetrics{}, err
	}

	var results []ServiceMetrics

	/*Loop and find the details of the listed services*/
	for _, s := range services.Services {

		for _, p := range processes {

			pName, _ := p.Name()
			if pName == s {
				cpu, _ := p.CPUPercent()
				status, _ := p.IsRunning()
				mem, _ := p.MemoryPercent()

				var statusResult string

				if status {
					statusResult = "Up"
				} else {
					statusResult = "Down"
				}

				results = append(results, ServiceMetrics{
					Name:   pName,
					Status: statusResult,
					CPU:    cpu,
					Memory: mem,
				})

				break
			}

		}
	}

	return results, nil
}
