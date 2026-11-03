package main

import (
 "context"

 "encoding/json"
 "log"
 "net/http"
 "net"
 "net/url"
 "os"
 "strings"
 "time"
)

type capabilities struct {
 Name string `json:"name"`
 GridConnected bool `json:"gridConnected"`
 ProtocolVersion string `json:"protocolVersion"`
 LoginSupported bool `json:"loginSupported"`
 RealtimeSupported bool `json:"realtimeSupported"`
 GridName string `json:"gridName"`
 LoginEndpointConfigured bool `json:"loginEndpointConfigured"`
}

func jsonResponse(w http.ResponseWriter, status int, value any) {
 w.Header().Set("Content-Type", "application/json; charset=utf-8")
 w.Header().Set("Cache-Control", "no-store")
 w.Header().Set("X-Content-Type-Options", "nosniff")
 w.WriteHeader(status)
 _ = json.NewEncoder(w).Encode(value)
}

// HTTP is supported for OpenSim G1, but credentials require an explicitly
// configured private transport. A public HTTP endpoint is never eligible.
func secureG1Transport(endpoint *url.URL) bool {
 if endpoint == nil { return false }
 if endpoint.Scheme == "https" { return true }
 if endpoint.Scheme != "http" || os.Getenv("G1_PRIVATE_TRANSPORT") != "true" { return false }
 host := endpoint.Hostname()
 ip := net.ParseIP(host)
 if ip == nil || !(ip.IsLoopback() || ip.IsPrivate()) { return false }
 return true
}

func main() {
 raw := strings.TrimSpace(os.Getenv("G1_LOGIN_URL"))
 var endpoint *url.URL
 if raw != "" {
  parsed, err := url.Parse(raw)
  if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") || parsed.User != nil {
   log.Fatal("G1_LOGIN_URL must be a valid http(s) URL without credentials")
  }
  endpoint = parsed
 }
 mux := http.NewServeMux()
 mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
  jsonResponse(w, http.StatusOK, map[string]string{"status":"ok"})
 })
 mux.HandleFunc("GET /api/v1/capabilities", func(w http.ResponseWriter, r *http.Request) {
  reachable := false
  if endpoint != nil {
   ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
   defer cancel()
   req, err := http.NewRequestWithContext(ctx, http.MethodHead, endpoint.String(), nil)
   if err == nil {
    client := &http.Client{Timeout: 3*time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
    resp, err := client.Do(req)
    if err == nil { reachable = true; resp.Body.Close() }
   }
  }
  jsonResponse(w, http.StatusOK, capabilities{
   Name:"Frivolin G1 WebGateway", GridName:"Stadt NexVerse – Generation 1",
   GridConnected:reachable, ProtocolVersion:"0.2",
   LoginSupported:secureG1Transport(endpoint) && os.Getenv("G1_ENABLE_LOGIN")=="true", RealtimeSupported:false, LoginEndpointConfigured:endpoint!=nil,
  })
 })
 mux.HandleFunc("POST /api/v1/session/login", func(w http.ResponseWriter, r *http.Request) {
  if !secureG1Transport(endpoint) || os.Getenv("G1_ENABLE_LOGIN") != "true" {
   jsonResponse(w, http.StatusServiceUnavailable, map[string]string{"code":"login_disabled","message":"Login gesperrt: HTTPS oder ein explizit freigegebener privater HTTP-Transport ist erforderlich."})
   return
  }
  handleG1Login(w, r, endpoint)

 })
 addr := os.Getenv("LISTEN_ADDR")
 if addr == "" { addr = "127.0.0.1:8090" }
 server := &http.Server{Addr:addr, Handler:mux, ReadHeaderTimeout:5*time.Second, ReadTimeout:10*time.Second, WriteTimeout:10*time.Second, IdleTimeout:30*time.Second}
 log.Printf("Frivolin G1 gateway listening on %s (G1 endpoint configured: %t)",addr,endpoint!=nil)
 log.Fatal(server.ListenAndServe())
}
