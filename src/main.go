package main

import (
	"fmt"
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

	go startServer(&wg, gConfigs, grpc_application.Initialize(gConfigs))
	go startServer(&wg, gConfigs, api.Initialize(gConfigs))

	time.Sleep(time.Duration(initDelay) * time.Millisecond)
	gConfigs.Display()

	// wait until WaitGroup is done
	wg.Wait()
	logger.Info("Good bye...")
}

func startServer(
	wg *sync.WaitGroup,
	gConfigs configs.GlobalConfigs,
	server network.ServerInterface,
) {
	logger.Info(fmt.Sprintf("Starting %s server...", server.GetName()))
	logger.Error(server.Run(gConfigs))
	wg.Done()
	logger.Warn(fmt.Sprintf("Server %s is gone...", server.GetName()))
}
