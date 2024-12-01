package mapper

import "github.com/go-clarum/agent/interface/grpc/agent/internal/api"

func TranslateCommand(command *api.ActionCommand) any {
	return nil
}

func TranslateResult(result any) *api.CommandResponse {
	return nil
}
