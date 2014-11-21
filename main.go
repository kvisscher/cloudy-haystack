package main

import (
    "flag"
    "net/http"
    "strconv"
    "log"
    "os"
    "github.com/kvisscher/cloudy-haystack/transform"
  )

func main() {
  var targetHost string
  var bindPort int

  flag.StringVar(&targetHost, "target", "http://somehost.com", "The host that the message should be re-posted to as JSON")
  flag.IntVar(&bindPort, "port", 8080, "The port to bind on to receive requests")
  flag.Parse()

  logFile, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)

  if err != nil {
    log.Fatal("Was unable to open log file", err)
  }

  defer logFile.Close()

  log.SetOutput(logFile)

  http.HandleFunc("/update", transform.UpdateHandler)
  http.ListenAndServe(":" + strconv.Itoa(bindPort), nil)
}
