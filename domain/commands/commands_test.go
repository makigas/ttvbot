package commands

import (
	"sync"
	"testing"
	"time"
)

func commandHandlerFactory(text string) commandHandler {
	return func(ctx *CommandContext) string { return text }
}

func TestGetCommands(t *testing.T) {
	list := NewCommandList()
	list.AddHandler("settitle", "change the stream title", commandHandlerFactory(""))
	list.AddHandler("debug", "test debug", commandHandlerFactory(""))
	cmds := list.Commands()

	if l := len(cmds); l != 2 {
		t.Errorf("Expected Commands() to retrieve 2 commands, retrieved %d", l)
	}
	if v, ok := cmds["settitle"]; !ok || v != "change the stream title" {
		t.Errorf("Expected settitle to be in the map, but it wasn't")
	}
	if v, ok := cmds["debug"]; !ok || v != "test debug" {
		t.Errorf("Expected debug to be in the map, but it wasn't")
	}
}

func TestExecuteCommandReturnsResult(t *testing.T) {
	list := NewCommandList()
	list.AddHandler("settitle", "change stream title", commandHandlerFactory("done"))

	ctx := CommandContext{}
	if _, err := list.Execute("notfound", &ctx); err == nil {
		t.Errorf("Expected Execute(notfound) to fail, but it didn't")
	}

	result, err := list.Execute("settitle", &ctx)
	if err != nil {
		t.Errorf("Expected Execute(settitle) not to fail, but it did: %s", err)
	}
	if result != "done" {
		t.Errorf("Expected function to yield `done`, but it did yield `%s`", result)
	}
}

func TestExecuteCommandExecutesCommand(t *testing.T) {
	var wg sync.WaitGroup
	callback := func(ctx *CommandContext) string {
		wg.Done()
		return "callback"
	}
	wg.Add(1)

	list := NewCommandList()
	ctx := CommandContext{}
	list.AddHandler("exec", "execute the handler", callback)
	if _, err := list.Execute("exec", &ctx); err != nil {
		t.Errorf("Expected Execute(exec) not to fail, but it did")
	}

	ch := waitGroupToChan(&wg)
	select {
	case <-ch:
		break
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestExecuteFallback(t *testing.T) {
	list := NewCommandList()
	list.AddHandler("test", "command that exists", commandHandlerFactory("hey"))
	list.SetFallback(commandHandlerFactory("generic fallback"))

	res, err := list.Execute("whatever", &CommandContext{})
	if err != nil {
		t.Errorf("Expected Execute to not fail because there is a fallback")
	}
	if res != "generic fallback" {
		t.Errorf("Expected Execute to yield `generic fallback`, but yielded `%s`", res)
	}
}
