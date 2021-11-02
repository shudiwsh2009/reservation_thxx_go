package main

import (
	"flag"
	"github.com/shudiwsh2009/reservation_thxx_go/service"
	"github.com/shudiwsh2009/reservation_thxx_go/web"
	"log"
	"os"
	"os/signal"
	"runtime"
)

func main() {
	var webAddress, confPath string
	var isDebug, isStaging bool
	flag.StringVar(&webAddress, "web", ":9000", "Web address server listening on")
	flag.StringVar(&confPath, "conf", "deploy/thxx.conf", "Configuration file path for service")
	flag.BoolVar(&isDebug, "debug", false, "Debug mode")
	flag.BoolVar(&isStaging, "staging", true, "Staging server")
	flag.Parse()

	service.InitService(confPath, isStaging)
	server := web.NewServer(isDebug)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		s := <-c
		log.Printf("Got signal: %s", s)
		server.Cleanup()
	}()

	log.Fatal(server.ListenAndServe(webAddress))
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
