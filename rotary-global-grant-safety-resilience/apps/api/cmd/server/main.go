package main

import (
  "log"
  "os"

  "github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/httpapi"
)

func main() {
  addr := getenv("HTTP_ADDR", ":8080")
  r := httpapi.NewRouter()

  log.Printf("api listening on %s", addr)
  if err := r.Run(addr); err != nil {
    log.Fatalf("server error: %v", err)
  }
}

func getenv(k, def string) string {
  v := os.Getenv(k)
  if v == "" {
    return def
  }
  return v
}
