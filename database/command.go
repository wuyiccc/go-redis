package database

import "strings"

type command struct {
	executor ExecFunc
	arity    int
}

var cmdTable = make(map[string]*command)

func RegisterCommand(name string, executor ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		executor: executor,
		arity:    arity,
	}
}
