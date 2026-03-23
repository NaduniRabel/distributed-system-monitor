package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	db "github.com/NaduniRabel/distributed-system-monitor/HQ/db"
)

type ServerStatus struct {
	HostID   string    `json:"host_id"`
	LastSeen time.Time `json:"last_seen"`
	Status   string    `json:"status"`
}

func main() {
	/*Initialize database*/
	db.Init()

	/*Route*/
	http.HandleFunc("/servers", getServersHandler)

	port := ":8080"

	server := &http.Server{
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("REST API  listening on %s...", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}

/*Get current hosts and statuses*/
func getServersHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	/*Find distinct host and calculate the status*/
	query := `
		SELECT DISTINCT ON (host_id) 
			host_id, 
			recorded_at,
			CASE 
				WHEN recorded_at > NOW() - INTERVAL '1 minute' THEN 'Online' 
				ELSE 'Offline' 
			END as status
		FROM system_logs 
		ORDER BY host_id, recorded_at DESC;`

	/*Execute query*/
	rows, err := db.Pool.Query(r.Context(), query)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var servers []ServerStatus
	for rows.Next() {
		var s ServerStatus
		if err := rows.Scan(&s.HostID, &s.LastSeen, &s.Status); err != nil {
			continue
		}
		servers = append(servers, s)
	}

	/*Return JSON*/
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}
