package main

import (
	"github.com/lgrossi/go-scaffold/src/api"
	"github.com/lgrossi/go-scaffold/src/configs"
	grpc_application "github.com/lgrossi/go-scaffold/src/grpc"
	"github.com/lgrossi/go-scaffold/src/logger"
	"github.com/lgrossi/go-scaffold/src/network"
	"sync"
	"time"
)

var numberOfServers = 2
var initDelay = 200

func main() {
	logger.Init(configs.GetLogLevel())
	logger.Info("Welcome to Your Fancy Application")
	logger.Info("Loading configurations...")

	var wg sync.WaitGroup
	wg.Add(numberOfServers)

	err := configs.Init()
	if err == nil {
		logger.Debug("Environment variables loaded from '.env'.")
	}

	gConfigs := configs.GetGlobalConfigs()

	go network.StartServer(&wg, gConfigs, &grpc_application.GrpcServer{})
	go network.StartServer(&wg, gConfigs, &api.Api{})

	time.Sleep(time.Duration(initDelay) * time.Millisecond)
	gConfigs.Display()

	// wait until WaitGroup is done
	wg.Wait()
	logger.Info("Good bye...")
}
