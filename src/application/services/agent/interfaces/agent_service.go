package interfaces

type AgentService interface {
	Status() string
	Shutdown()
}
