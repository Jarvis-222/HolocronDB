package core

import (
	"fmt"
	"log"
	"strconv"
)

var RES_NIL = []byte("$-1\r\n")

func EvalCommand(cmd *RedisCommand) (interface{}, error) {

	log.Printf("Evaluating command %s with args %v", cmd.Cmd, cmd.Args)

	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd)
	case "SET":
		return evalSET(cmd)
	case "GET":
		return evalGET(cmd)
	case "TTL":
		return evalTTL(cmd)
	default:
		return nil, fmt.Errorf("unknown command '%s'", cmd.Cmd)
	}

}

func evalPING(cmd *RedisCommand) (interface{}, error) {

	if len(cmd.Args) == 0 {
		return Encode("PONG", true), nil
	}

	if len(cmd.Args) >= 2 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", cmd.Cmd)
	}

	return Encode(cmd.Args[0], false), nil
}

func evalSET(cmd *RedisCommand) (interface{}, error) {

	args := cmd.Args

	if len(args) < 2 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", cmd.Cmd)
	}

	var key, val string
	var expDurationMs int64 = -1

	key = args[0]
	val = args[1]

	for i := 2; i < len(args); i++ {

		switch args[i] {
		case "EX", "ex":
			i++
			if i == len(args) {
				return nil, fmt.Errorf("syntax error")
			}
			expDurationSec, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("syntax error")
			}
			expDurationMs = expDurationSec * 1000

		default:
			return nil, fmt.Errorf("syntax error")
		}

	}
	Put(key, val, expDurationMs)
	return Encode("OK", false), nil

}

func evalGET(cmd *RedisCommand) (interface{}, error) {

	args := cmd.Args
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", cmd.Cmd)
	}

	key := args[0]
	val, exists := Get(key)

	if !exists {
		return Encode(nil, false), nil
	}

	return Encode(val, false), nil
}

func evalTTL(cmd *RedisCommand) (interface{}, error) {

	args := cmd.Args
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", cmd.Cmd)
	}

	key := args[0]
	ttl := GetTTL(key)
	return Encode(ttl, false), nil

}

func evalDelete(cmd *RedisCommand) (interface{}, error) {

	args := cmd.Args

	if len(args) < 1 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", cmd.Cmd)
	}

	cntDeleted := 0; 

	for i := range args {
		key := args[i]
		cntDeleted += Delete(key) 
	}
	
	return Encode(cntDeleted, false), nil
}

func evalExpire(cmd *RedisCommand) (interface{}, error) {

	args := cmd.Args

	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", cmd.Cmd)
	}
	
	key := args[0]
	expDurationSec, err := strconv.ParseInt(args[1], 10, 64)

	if err != nil {
		return nil, fmt.Errorf("syntax error")
	}

	cntUpdated := ExpireAt(key, expDurationSec * 1000)
	
	return Encode(cntUpdated, false), nil
}