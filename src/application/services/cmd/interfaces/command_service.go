package interfaces

type CommandService interface {
	InitializeEndpoint(name string, cmdComponents []string, warmupMillis int32) error
	ShutdownEndpoint(name string) error
}
