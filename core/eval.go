package core

import (
	"fmt"
	"log"
)

func EvalCommand(cmd *RedisCommand) (interface{}, error) {

	log.Printf("Evaluating command %s with args %v", cmd.Cmd, cmd.Args)

	switch cmd.Cmd {
	case "PING":
		return EvalPingCommand(cmd)
	default:
		return EvalPingCommand(cmd)
	}

}

func EvalPingCommand(cmd *RedisCommand) (interface{}, error) {

	if len(cmd.Args) == 0 {
		return Encode("PONG", true), nil
	}

	if len(cmd.Args) >= 2 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", cmd.Cmd)
	}

	return Encode(cmd.Args[0], false), nil
}
