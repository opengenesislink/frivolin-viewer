package main

import (
  "encoding/json"
  "log"
  "net/http"
  "os"
  "time"
)

type capabilities struct {
  Name string `json:"name"`
  GridConnected bool `json:"gridConnected"`
  ProtocolVersion string `json:"protocolVersion"`
  LoginSupported bool `json:"loginSupported"`
  RealtimeSupported bool `json:"realtimeSupported"`
}

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
  })
  mux.HandleFunc("GET /api/v1/capabilities", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Cache-Control", "no-store")
    _ = json.NewEncoder(w).Encode(capabilities{
      Name: "Frivolin G1 WebGateway scaffold",
      GridConnected: false, ProtocolVersion: "0.1",
      LoginSupported: false, RealtimeSupported: false,
    })
  })
  addr := os.Getenv("LISTEN_ADDR")
  if addr == "" { addr = "127.0.0.1:8090" }
  server := &http.Server{
    Addr: addr,
    Handler: mux,
    ReadHeaderTimeout: 5 * time.Second,
    IdleTimeout: 30 * time.Second,
  }
  log.Printf("Frivolin gateway scaffold listening on %s", addr)
  log.Fatal(server.ListenAndServe())
}
