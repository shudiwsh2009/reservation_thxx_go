package main

import (
	"github.com/shudiwsh2009/reservation_thxx_go/service"
	"github.com/shudiwsh2009/reservation_thxx_go/web"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
)

func main() {
	var webAddress, debugAssetsPort, confPath string
	var isDebug, isStaging bool
	flag.StringVar(&webAddress, "web", ":9000", "Web address server listening on")
	flag.StringVar(&debugAssetsPort, "devWeb", "", "Web address server listening on (like :9010)")
	flag.StringVar(&confPath, "conf", "deploy/thxx.conf", "Configuration file path for service")
	flag.BoolVar(&isDebug, "debug", false, "Debug mode")
	flag.BoolVar(&isStaging, "staging", true, "Staging server")
	flag.Parse()

	service.InitService(confPath, isStaging)
	server := web.NewServer(isDebug)
	if isDebug && debugAssetsPort != "" {
		server.SetAssetDomain("//localhost" + debugAssetsPort)
	}
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
