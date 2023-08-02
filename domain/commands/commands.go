package commands

import "fmt"

// CommandContext provides information about the executed command.
type CommandContext struct {
	Username   string   // username that invoked the command
	Vip        bool     // if the username is vip
	Subscriber bool     // if the username is subscriber
	Mod        bool     // if the username is mod
	Streamer   bool     // if the username is streamer
	Params     []string // the list of params
	Message    string   // the complete full message
}

type commandHandler func(ctx *CommandContext) string

type command struct {
	description string
	handler     commandHandler
}

type CommandList struct {
	handlers map[string]command
	fallback commandHandler
}

func NewCommandList() *CommandList {
	cmd := CommandList{
		handlers: make(map[string]command),
		fallback: nil,
	}
	return &cmd
}

func (list *CommandList) AddHandler(name, description string, callback commandHandler) {
	list.handlers[name] = command{
		description: description,
		handler:     callback,
	}
}

func (list *CommandList) SetFallback(callback commandHandler) {
	list.fallback = callback
}

func (list *CommandList) Commands() map[string]string {
	data := make(map[string]string)
	for k, v := range list.handlers {
		data[k] = v.description
	}
	return data
}

func (list *CommandList) Execute(name string, ctx *CommandContext) (string, error) {
	if command, ok := list.handlers[name]; ok {
		return command.handler(ctx), nil
	}
	if list.fallback != nil {
		return list.fallback(ctx), nil
	}
	return "", fmt.Errorf("command %s not found", name)
}
