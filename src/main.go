package main

import (
	"fmt"
	"github.com/go-clarum/agent/infrastructure/config"
	"github.com/go-clarum/agent/infrastructure/logging"
	"github.com/go-clarum/agent/interface/grpc/agent"
	"google.golang.org/grpc"
	"net"
)

func main() {
	logging.Infof("starting clarum agent v%s", config.Version())

	initAndRunGrpcServer()

	logging.Info("shutting down clarum agent")
}

func initAndRunGrpcServer() {
	address := fmt.Sprintf("localhost:%d", config.AgentPort())
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logging.Errorf("failed to initiate GRPC server on port [%d]: %s", config.AgentPort(), err)
		return
	}

	grpcServer := grpc.NewServer()

	agent.RegisterAgentApi(grpcServer)

	logging.Infof("starting GRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		logging.Errorf("GRPC server startup error: %s", err)
	}
}
