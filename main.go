package main

import (
	"fmt"
	"github.com/go-clarum/agent/config"
	"github.com/go-clarum/agent/control"
	"github.com/go-clarum/agent/logging"
	"github.com/go-clarum/agent/services/agent"
	"github.com/go-clarum/agent/services/cmd"
	"github.com/go-clarum/agent/services/http"
	"google.golang.org/grpc"
	"net"
)

func main() {
	logging.Infof("Starting clarum agent v%s", config.Version())
	control.ShutdownHook.Add(1)

	config.LoadConfig()
	initAndRunGrpcServer()

	control.ShutdownHook.Wait()
	logging.Info("Shutting down clarum agent")
}

func initAndRunGrpcServer() {
	address := fmt.Sprintf("localhost:%d", config.AgentPort())
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logging.Errorf("Failed to initiate GRPC server on port [%d]: %s", config.AgentPort(), err)
		return
	}

	grpcServer := grpc.NewServer()

	agent.RegisterAgentService(grpcServer)
	http.RegisterHttpService(grpcServer)
	cmd.RegisterCmdService(grpcServer)

	logging.Infof("Starting GRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		return
	}
}
