package server

import (
	"HolocronDB/config"
	"HolocronDB/core"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

func DecodeStringArray(data []byte) ([]string, error) {

	arr, err := core.Decode(data)

	if err != nil {
		return nil, err
	}

	ts := arr.([]any)
	tokens := make([]string, len(ts))

	for i := range tokens {
		tokens[i] = ts[i].(string)
	}

	return tokens, nil
}

func readCommand(c io.ReadWriter) (*core.RedisCommand, error) {

	var buf []byte = make([]byte, 512)

	n, err := c.Read(buf[:])

	if err != nil {
		return nil, err
	}

	args, err := DecodeStringArray(buf[:n])

	if err != nil {
		return nil, err
	}

	for i := range args {
		log.Println(args[i])
	}
	return &core.RedisCommand{
		Cmd:  args[0],
		Args: args[1:],
	}, nil

}

func respond(cmd *core.RedisCommand, c io.ReadWriter) {

	res, err := core.EvalCommand(cmd)

	if err != nil {
		log.Println("Error occured: ", err)
		respondError(err, c)
	} else {
		c.Write(res.([]byte))
	}

}

func respondError(err error, c io.ReadWriter) {

	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func RunSyncTCPServer() {

	log.Println("Starting synchronous TCP server...", config.Host, config.Port)

	var con_clients int = 0

	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))

	if err != nil {
		panic(err)
	}

	for {

		c, err := lsnr.Accept()

		if err != nil {
			panic(err)
		}

		con_clients += 1

		log.Println("New client connected with address: ", c.RemoteAddr(), " | con-current clients: ", con_clients)

		for {

			cmd, err := readCommand(c)

			if err != nil {
				c.Close()
				log.Println("Client disconnected with address: ", c.RemoteAddr(), " | con-current clients: ", con_clients-1)
				con_clients -= 1

				if err.Error() == "EOF" {
					break
				}
				log.Println("Error: ", err)

			}
			log.Println("Received command: ", cmd)

			respond(cmd, c)

		}

	}
}
