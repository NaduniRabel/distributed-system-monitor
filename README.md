<h1>Distributed System Monitor</h1>

This project is a real-time monitoring system. It uses a Go-based Agent to collect metrics and a gRPC HQ Server to store them in PostgreSQL. 
A REST API is provided to view server details.

<h2>Prerequisites</h2>

-Go 1.21+ installed.
-Protobuf installed
-PostgreSQL installed and running.
-Git to clone the repository.

<h2>Database Setup</h2>
Before starting the services, ensure your PostgreSQL instance is running.
Create a database named with a preferred name. 

```bash
// Format: postgres://<user>:<password>@localhost:5432/<database_name>?sslmode=disable
databaseURL := "postgres://user:password@localhost:5432/postgres?sslmode=disable"
```

<h2>Running the HQ (gRPC Server & REST API)</h2>

The HQ consists of two main components. Open two separate terminals to run them:

1. Start the gRPC Server (Collector)
This server listens on port 9001 for incoming streams from agents.

```bash
cd HQ/grpcServer
go mod tidy
go run main.go
```

2. Start the REST API
This server listens on port 8080 to serve server status data.

```bash
cd HQ/restApi
go mod tidy
go run main.go
```
Access via:
```bash
http://localhost:8080
```

<h2>Running the Agent</h2>
The Agent collects local CPU, Memory, and Disk metrics and streams them to the HQ. Open a new terminal:
Ensure the gRPC server address in agent/main.go matches your HQ location (default is localhost:9001).

```bash
cd agent
go mod tidy
go run main.go
```




