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
	logging.Infof("starting clarum agent v%s", config.Version())
	control.ShutdownHook.Add(1)

	initAndRunGrpcServer()

	control.ShutdownHook.Wait()
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

	agent.RegisterAgentService(grpcServer)
	http.RegisterHttpService(grpcServer)
	cmd.RegisterCmdService(grpcServer)

	logging.Infof("starting GRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		logging.Errorf("GRPC server startup error: %s", err)
	}
}
