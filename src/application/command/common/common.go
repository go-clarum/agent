package common

type CommandHandler interface {
	CanHandle(command any) bool
	Handle(command any) any
}
